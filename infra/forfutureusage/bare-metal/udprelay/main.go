// A UDP port-range relay meant to run on the WINDOWS host (not inside WSL2).
// It listens on every port in [start, end] on the Windows side and forwards
// packets to the same port on your minikube/WSL2 IP, then relays replies
// back to whichever client last talked to that port.
//
// GOOS=windows GOARCH=amd64 go build -o <filename>.exe
//
// Run (as Administrator):
//   <filename>.exe -dest 172.x.x.x -start 50000 -end 60000

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	readBufSize = 64 * 1024       // generous; RTP/STUN packets are small but be safe
	idleTimeout = 5 * time.Minute // Time before a UDP session is considered stale
	sweepPeriod = 1 * time.Minute // How often the garbage collector runs
)

// portRelay holds the state for a single relayed UDP port.
type portRelay struct {
	port       int
	listenConn *net.UDPConn // socket clients on the internet talk to
	destConn   *net.UDPConn // socket dialed to $MINIKUBE_IP:port

	mu         sync.Mutex
	lastClient *net.UDPAddr // most recent external client address seen
	lastSeen   time.Time
}

func newPortRelay(listenIP string, destIP string, port int) (*portRelay, error) {
	laddr := &net.UDPAddr{IP: net.ParseIP(listenIP), Port: port}
	lconn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, fmt.Errorf("listen :%d: %w", port, err)
	}

	daddr := &net.UDPAddr{IP: net.ParseIP(destIP), Port: port}
	dconn, err := net.DialUDP("udp4", nil, daddr)
	if err != nil {
		lconn.Close()
		return nil, fmt.Errorf("dial dest :%d: %w", port, err)
	}

	pr := &portRelay{
		port:       port,
		listenConn: lconn,
		destConn:   dconn,
	}
	return pr, nil
}

func (pr *portRelay) close() {
	pr.listenConn.Close()
	pr.destConn.Close()
}

// runClientToDest reads packets arriving from the internet and forwards
// them to the destination (minikube), remembering who sent them.
func (pr *portRelay) runClientToDest(wg *sync.WaitGroup, errCount *int64) {
	defer wg.Done()
	buf := make([]byte, readBufSize)
	for {
		n, raddr, err := pr.listenConn.ReadFromUDP(buf)
		if err != nil {
			if isClosedErr(err) {
				return
			}
			atomic.AddInt64(errCount, 1)
			continue
		}

		pr.mu.Lock()
		pr.lastClient = raddr
		pr.lastSeen = time.Now()
		pr.mu.Unlock()

		if _, err = pr.destConn.Write(buf[:n]); err != nil {
			atomic.AddInt64(errCount, 1)
		}
	}
}

// runDestToClient reads replies coming back from minikube and forwards
// them to whichever client last sent a packet on this port.
func (pr *portRelay) runDestToClient(wg *sync.WaitGroup, errCount *int64) {
	defer wg.Done()
	buf := make([]byte, readBufSize)
	for {
		n, err := pr.destConn.Read(buf)
		if err != nil {
			if isClosedErr(err) {
				return
			}
			atomic.AddInt64(errCount, 1)
			continue
		}

		pr.mu.Lock()
		client := pr.lastClient
		pr.mu.Unlock()

		if client == nil {
			// No client has talked to us yet on this port; drop.
			continue
		}
		if _, err = pr.listenConn.WriteToUDP(buf[:n], client); err != nil {
			atomic.AddInt64(errCount, 1)
		}
	}
}

// isClosedErr reports whether err is the result of the socket being closed
func isClosedErr(err error) bool {
	return err != nil && errors.Is(err, net.ErrClosed)
}

// startSweeper runs in the background and clears out stale client addresses
func startSweeper(relays []*portRelay) {
	ticker := time.NewTicker(sweepPeriod)
	go func() {
		for range ticker.C {
			now := time.Now()
			clearedCount := 0

			for _, pr := range relays {
				pr.mu.Lock()
				// If we have a client AND it has been silent longer than our timeout threshold
				if pr.lastClient != nil && now.Sub(pr.lastSeen) > idleTimeout {
					pr.lastClient = nil
					clearedCount++
				}
				pr.mu.Unlock()
			}

			if clearedCount > 0 {
				log.Printf("sweeper: cleared %d stale UDP allocations", clearedCount)
			}
		}
	}()
}

func main() {
	var (
		destIP   = flag.String("dest", "", "destination IP to forward to, e.g. your minikube IP (required)")
		listenIP = flag.String("listen", "0.0.0.0", "IP to listen on (default: all interfaces)")
		start    = flag.Int("start", 0, "first port in range (required)")
		end      = flag.Int("end", 0, "last port in range, inclusive (required)")
	)
	flag.Parse()

	if *destIP == "" || *start == 0 || *end == 0 || *end < *start {
		fmt.Fprintln(os.Stderr, "usage: udprelay -dest <minikube-ip> -start <port> -end <port> [-listen <ip>]")
		os.Exit(1)
	}

	portCount := *end - *start + 1
	log.Printf("starting UDP relay: %s:[%d-%d] -> %s:[%d-%d] (%d ports)",
		*listenIP, *start, *end, *destIP, *start, *end, portCount)

	var wg sync.WaitGroup
	var errCount int64
	relays := make([]*portRelay, 0, portCount)

	failed := 0
	for p := *start; p <= *end; p++ {
		pr, err := newPortRelay(*listenIP, *destIP, p)
		if err != nil {
			failed++
			if failed <= 10 {
				log.Printf("WARNING: could not set up port %d: %v", p, err)
			} else if failed == 11 {
				log.Printf("WARNING: suppressing further per-port setup errors...")
			}
			continue
		}
		relays = append(relays, pr)

		wg.Add(2)
		go pr.runClientToDest(&wg, &errCount)
		go pr.runDestToClient(&wg, &errCount)

		if p%1000 == 0 {
			log.Printf("...bound through port %d", p)
		}
	}

	log.Printf("relay active: %d/%d ports bound successfully (%d failed)",
		len(relays), portCount, failed)
	if len(relays) == 0 {
		log.Fatal("no ports bound successfully, exiting")
	}

	// Kick off the background garbage collector
	startSweeper(relays)

	// Graceful shutdown on Ctrl+C / service stop.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("shutting down, closing all sockets...")
	for _, pr := range relays {
		pr.close()
	}
	wg.Wait()
	log.Printf("done. total forwarding errors during run: %d", atomic.LoadInt64(&errCount))
}

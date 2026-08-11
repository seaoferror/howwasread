package manager

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pion/turn/v5"
)

type TurnManager struct {
	packetConn net.PacketConn
	mu         sync.Mutex
	turnServer *turn.Server
	publicIP   net.IP
}

func RunTurnManager(packetConn net.PacketConn, checkInterval time.Duration) (*TurnManager, error) {
	tm := &TurnManager{
		packetConn: packetConn,
	}

	initialIP, err := fetchPublicWANIP()
	if err != nil {
		return nil, fmt.Errorf("could not determine initial WAN IP: %w", err)
	}

	err = tm.serve(initialIP)
	if err != nil {
		return nil, err
	}

	go tm.monitorAndReflectIPChange(checkInterval)

	return tm, nil
}

func (tm *TurnManager) serve(publicIP net.IP) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.turnServer != nil {
		log.Printf("[TURN] Public IP changed to %s. Closing old TURN server...", publicIP.String())
		err := tm.turnServer.Close()
		if err != nil {
			log.Printf("[TURN] Error closing old server instance: %v", err)
		}
	}

	log.Printf("[TURN] Starting server instance bound to RelayAddress: %s", publicIP.String())
	server, err := turn.NewServer(turn.ServerConfig{
		Realm: os.Getenv("TURN_REALM"),
		AuthHandler: func(ra *turn.RequestAttributes) (string, []byte, bool) {
			username := ra.Username
			parts := strings.Split(username, ":")
			timestampStr := parts[0]
			expiry, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil || time.Now().Unix() > expiry {
				return "", nil, false
			}
			mac := hmac.New(sha1.New, []byte(os.Getenv("TURN_SECRET")))
			mac.Write([]byte(username))
			expectedPassword := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			return username, turn.GenerateAuthKey(username, ra.Realm, expectedPassword), true
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: tm.packetConn,
				RelayAddressGenerator: &turn.RelayAddressGeneratorPortRange{
					RelayAddress: publicIP,
					MinPort:      42093,
					MaxPort:      49000,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create turn server: %w", err)
	}

	tm.turnServer = server
	tm.publicIP = publicIP
	return nil
}

func (tm *TurnManager) Close() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.turnServer != nil {
		return tm.turnServer.Close()
	}
	return nil
}

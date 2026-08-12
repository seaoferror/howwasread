package manager

import (
	"backend/common"
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
	realm      string
	secret     string
	publicIP   net.IP
	zoneID     string
	recordID   string
	apiToken   string
}

func RunTurnManager(packetConn net.PacketConn, checkInterval time.Duration) (*TurnManager, error) {
	tm := &TurnManager{
		packetConn: packetConn,
		realm:      os.Getenv("TURN_REALM"),
		secret:     os.Getenv("TURN_SECRET"),
		zoneID:     os.Getenv("CF_ZONE_ID"),
		recordID:   os.Getenv("CF_RECORD_ID"),
		apiToken:   os.Getenv("CF_API_TOKEN"),
	}
	if tm.realm == "" || tm.secret == "" || tm.zoneID == "" || tm.recordID == "" || tm.apiToken == "" {
		return nil, fmt.Errorf("cloudflare credentials missing in environment variables")
	}
	initialIP, err := common.FetchPublicWANIP()
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
		newConn, err := net.ListenPacket("udp4", "0.0.0.0:3478")
		if err != nil {
			return fmt.Errorf("failed to bind UDP port 3478: %w", err)
		}
		tm.packetConn = newConn
	}

	log.Printf("[TURN] Starting server instance bound to RelayAddress: %s", publicIP.String())
	server, err := turn.NewServer(turn.ServerConfig{
		Realm: tm.realm,
		AuthHandler: func(ra *turn.RequestAttributes) (string, []byte, bool) {
			username := ra.Username
			parts := strings.Split(username, ":")
			timestampStr := parts[0]
			expiry, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil || time.Now().Unix() > expiry {
				return "", nil, false
			}
			mac := hmac.New(sha1.New, []byte(tm.secret))
			mac.Write([]byte(username))
			credential := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			return username, turn.GenerateAuthKey(username, ra.Realm, credential), true
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: tm.packetConn,
				RelayAddressGenerator: &turn.RelayAddressGeneratorPortRange{
					RelayAddress: publicIP,
					Address:      "0.0.0.0",
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

package manager

import (
	"backend/common"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
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
	udp        net.PacketConn
	tcp        net.Listener
	tlsConfig  *tls.Config
	mu         sync.Mutex
	turnServer *turn.Server
	realm      string
	secret     string
	publicIP   net.IP
	zoneID     string
	recordID   string
	apiToken   string
}

func RunTurnManager(checkInterval time.Duration) (tm *TurnManager, err error) {
	tm = &TurnManager{
		realm:    os.Getenv("TURN_REALM"),
		secret:   os.Getenv("TURN_SECRET"),
		zoneID:   os.Getenv("CF_ZONE_ID"),
		recordID: os.Getenv("CF_TURN_RECORD_ID"),
		apiToken: os.Getenv("CF_API_TOKEN"),
	}
	tm.tlsConfig, err = common.CreateTlSConfig(os.Getenv("TURN_CERT_PATH"), os.Getenv("TURN_KEY_PATH"), "")
	if err != nil {
		return nil, err
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
			return err
		}
		u, err := net.ListenPacket("udp4", "0.0.0.0:3478")
		if err != nil {
			return err
		}
		tm.udp = u
		t, err := tls.Listen("tcp", "0.0.0.0:5349", tm.tlsConfig)
		if err != nil {
			return err
		}
		tm.tcp = t
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
		ListenerConfigs: []turn.ListenerConfig{
			{
				Listener: tm.tcp,
				RelayAddressGenerator: &turn.RelayAddressGeneratorPortRange{
					RelayAddress: publicIP,
					Address:      "0.0.0.0",
					MinPort:      42093,
					MaxPort:      49000,
				},
			},
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: tm.udp,
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

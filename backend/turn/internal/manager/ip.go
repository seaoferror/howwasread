package manager

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func fetchPublicWANIP() (net.IP, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public IP: %w", err)
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ipStr := strings.TrimSpace(string(ipBytes))
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid IP format received: %s", ipStr)
	}

	return parsedIP, nil
}

func (tm *TurnManager) monitorAndReflectIPChange(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		latestIP, err := fetchPublicWANIP()
		if err != nil {
			log.Printf("[TURN Monitor] Warning: failed to check WAN IP: %v", err)
			continue
		}

		tm.mu.Lock()
		activeIP := tm.publicIP
		tm.mu.Unlock()

		if !latestIP.Equal(activeIP) {
			log.Printf("[TURN Monitor] IP change detected! Old: %s -> New: %s", activeIP.String(), latestIP.String())
			err = tm.serve(latestIP)
			if err != nil {
				log.Printf("[TURN Monitor] Critical: Failed to restart TURN server with new IP: %v", err)
			}
		}
	}
}

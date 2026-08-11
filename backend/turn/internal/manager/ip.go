package manager

import (
	"backend/common"
	"log"
	"time"
)

func (tm *TurnManager) monitorAndReflectIPChange(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		latestIP, err := common.FetchPublicWANIP()
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

package profile

import (
	"github.com/0studio/scheduler"
	log "github.com/cihub/seelog"
	"time"
)

const MIN_CHECK_SECONDS = 60 // 定时器可接受的最小间隔
func startCleaner(m *LocalUserStorage) {
	timer := scheduler.NewTimingWheel(MIN_CHECK_SECONDS*time.Second, 2)
	go func() {
		for {
			select {
			case <-timer.After(MIN_CHECK_SECONDS * time.Second):
				log.Debug("trying clean cacheuser")
				for key, cacheObj := range m.cache {
					_, ok := cacheObj.GetObject(time.Now())
					if !ok {
						m.cache.Delete(key)
						log.Debug("delete outdate cacheuser", key)
					}
				}
			}
		}
	}()
}

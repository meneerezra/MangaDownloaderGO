package fetcher

import (
	"mangaDownloaderGO/utils/logger"
	"time"
)

type RateLimit struct {
	TimeoutSeconds time.Duration
	TimeLastUsed time.Time
}

func (rateLimit *RateLimit) HandleRatelimit() {
	rateLimit.TimeoutSeconds += 1
	rateLimitTimer := time.NewTimer(rateLimit.TimeoutSeconds * time.Second)
	logger.WarningFromStringF("Rate limit hit starting %v timer", rateLimit.TimeoutSeconds * time.Second)
	<-rateLimitTimer.C
	rateLimit.TimeLastUsed = time.Now()
}

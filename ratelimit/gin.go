package ratelmit

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	maxRequests     = 2
	perMinutePeriod = 1 * time.Minute
)

var (
	ipRequestsCounts = make(map[string]int) //Can use some distributed db
	mutex            = &sync.Mutex{}
)

func RateLimiter(context *gin.Context) {
	ip := context.ClientIP()
	mutex.Lock()
	defer mutex.Unlock()
	count := ipRequestsCounts[ip]
	if count >= maxRequests {
		context.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	ipRequestsCounts[ip] = count + 1
	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
	})

	context.Next()
}

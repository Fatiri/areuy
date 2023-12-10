package authentication

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Fatiri/areuy/exception"
	"github.com/gin-gonic/gin"
)

const (
	maxRequests     = 100
	perMinutePeriod = 1 * time.Minute
)

var (
	ipRequestsCounts = make(map[string]int) //Can use some distributed db
	mutex            = &sync.Mutex{}
)

func RateLimiterGin(access string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		mutex.Lock()
		defer mutex.Unlock()
		count := ipRequestsCounts[ip]
		if count >= maxRequests {
			ctx.JSON(http.StatusForbidden, exception.Error(errors.New("there are too many requests"), exception.Message{
				Id: "Permintaannya terlalu banyak, silahkan coba lagi dalam 1 menit!",
				En: "There are too many requests, please try again in 1 minute!",
			}, access))
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		if len(ipRequestsCounts) >= 1000 {
			ctx.JSON(http.StatusForbidden, exception.Error(errors.New("there are too many requests"), exception.Message{
				Id: "Permintaannya terlalu banyak, silahkan coba lagi dalam 1 menit!",
				En: "There are too many requests, please try again in 1 minute!",
			}, access))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ipRequestsCounts[ip] = count + 1
		time.AfterFunc(perMinutePeriod, func() {
			mutex.Lock()
			defer mutex.Unlock()

			ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
		})

		ctx.Next()
	}
}

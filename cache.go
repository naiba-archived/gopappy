package gopappy

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Ch *cache.Cache

func init() {
	if Ch == nil {
		Ch = cache.New(time.Minute*20, time.Minute*40)
	}
}

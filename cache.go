// Copyright (c) 2018 奶爸(1@5.nu)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

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

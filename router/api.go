/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"git.cm/naiba/gopappy"
	"git.cm/naiba/gopappy/util/aliyun"
	"git.cm/naiba/gopappy/util/cn4"
	"git.cm/naiba/gopappy/util/ename"
	"github.com/gin-gonic/gin"
)

func serveAPI(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.POST("search", search)
		api.GET("params", params)
	}
}

func search(c *gin.Context) {
	var o gopappy.Option
	if err := c.BindJSON(&o); err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	var wg sync.WaitGroup
	funks := map[string]func(option gopappy.Option) ([]gopappy.Domain, error){
		"ename": ename.Domains,
		"cn4":   cn4.Domains,
		"ali":   aliyun.Domains,
	}
	type SafeDomains struct {
		D []gopappy.Domain
		L sync.RWMutex
	}
	var all SafeDomains
	all.D = make([]gopappy.Domain, 0)
	// genCacheKey
	o.Page--
	hasKey := genCacheKey(o)
	o.Page++
	cacheKey := genCacheKey(o)
	// range get domains
	for p, f := range funks {
		wg.Add(1)
		go func(p string, f func(option gopappy.Option) ([]gopappy.Domain, error)) {
			d := getCachedDomains(p+hasKey, p+cacheKey, o, f)
			all.L.Lock()
			all.D = append(all.D, d...)
			all.L.Unlock()
			wg.Done()
		}(p, f)
	}
	wg.Wait()
	convertPrices(all.D)
	if o.Sort == 2 {
		sort.Sort(sort.Reverse(gopappy.SortDomain(all.D)))
	} else {
		sort.Sort(gopappy.SortDomain(all.D))
	}
	gopappy.StatInstance.Search++
	gopappy.StatInstance.Domain += int64(len(all.D))
	c.JSON(http.StatusOK, all.D)
}

func params(c *gin.Context) {
	type TLD struct {
		ID  int    `json:"id,omitempty"`
		TLD string `json:"tld,omitempty"`
	}
	type Params struct {
		Platforms map[int]string `json:"platforms,omitempty"`
		TLDs      []TLD          `json:"tlds,omitempty"`
		Stat      *gopappy.Stat  `json:"stat,omitempty"`
	}
	var p = Params{
		Platforms: gopappy.Platform,
		Stat:      gopappy.StatInstance,
	}
	p.TLDs = make([]TLD, 0)
	for id, tld := range gopappy.TLDs {
		p.TLDs = append(p.TLDs, TLD{ID: id, TLD: tld})
	}
	c.JSON(http.StatusOK, p)
}

func getCachedDomains(hasKey, cacheKey string, o gopappy.Option, fn func(option gopappy.Option) ([]gopappy.Domain, error)) []gopappy.Domain {
	domains := make([]gopappy.Domain, 0)
	// 是否有下一页
	_, has := gopappy.Ch.Get("H" + hasKey)
	log.Println(hasKey, has)
	if has {
		return domains
	}
	// 取缓存或重新抓取
	tmp, has := gopappy.Ch.Get(cacheKey)
	log.Println(cacheKey, has)
	if has {
		return tmp.([]gopappy.Domain)
	} else {
		var err error
		domains, err = fn(o)
		if err == nil {
			gopappy.Ch.Set(cacheKey, domains, time.Minute*30)
			if len(domains) < 20 {
				gopappy.Ch.Set("H"+cacheKey, false, time.Minute*30)
			}
		}
		return domains
	}
}

func genCacheKey(o gopappy.Option) string {
	cacheKey := fmt.Sprintf("EX%sEXP%dKW%sKWP%dML%dMP%dIL%dIP%dOD%dS%dP%dT%dTD", o.Exclude, o.ExPos,
		o.Keyword, o.KwPos, o.MaxLength, o.MaxPrice, o.MinLength, o.MinPrice,
		o.Order, o.Sort, o.Page, o.Tag)
	sort.Sort(sort.IntSlice(o.TLDs))
	for _, tld := range o.TLDs {
		cacheKey += strconv.Itoa(tld) + ","
	}
	return cacheKey
}

func convertPrices(ds []gopappy.Domain) {
	for i, d := range ds {
		if d.Currency != "CNY" && len(d.Currency) > 0 {
			d.Price = int(gopappy.Fx.Convert(d.Currency, "CNY", float64(d.Price)))
			d.Currency = "CNY"
		}
		if d.Price < 1000 && d.Price > 0 {
			d.PriceString = fmt.Sprintf("%d 元", d.Price)
		} else if d.Price < 10000 {
			d.PriceString = fmt.Sprintf("%d 千 %d 元", d.Price/1000, d.Price%1000)
		} else {
			d.PriceString = fmt.Sprintf("%d 万 %d 千 %d 元", d.Price/10000, d.Price%10000/1000, d.Price%10000%1000)
		}
		ds[i] = d
	}
}

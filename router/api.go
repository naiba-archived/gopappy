/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

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
	funks := []func(option gopappy.Option) ([]gopappy.Domain, error){
		ename.Domains,
		cn4.Domains,
		aliyun.Domains,
	}
	type SafeDomains struct {
		D []gopappy.Domain
		L sync.RWMutex
	}
	var all SafeDomains
	all.D = make([]gopappy.Domain, 0)
	for _, f := range funks {
		wg.Add(1)
		go func(f func(option gopappy.Option) ([]gopappy.Domain, error)) {
			d, _ := f(o)
			all.L.Lock()
			all.D = append(all.D, d...)
			all.L.Unlock()
			wg.Done()
		}(f)
	}
	wg.Wait()
	convertPrices(all.D)
	if o.Sort == 2 {
		sort.Sort(sort.Reverse(gopappy.SortDomain(all.D)))
	} else {
		sort.Sort(gopappy.SortDomain(all.D))
	}
	c.JSON(http.StatusOK, all.D)
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

func params(c *gin.Context) {
	type TLD struct {
		ID  int    `json:"id"`
		TLD string `json:"tld"`
	}
	type Params struct {
		Platforms map[int]string `json:"platforms"`
		TLDs      []TLD          `json:"tlds"`
	}
	var p = Params{
		Platforms: gopappy.Platform,
	}
	p.TLDs = make([]TLD, 0)
	for id, tld := range gopappy.TLDs {
		p.TLDs = append(p.TLDs, TLD{ID: id, TLD: tld})
	}
	c.JSON(http.StatusOK, p)
}

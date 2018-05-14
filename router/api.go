/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"fmt"
	"git.cm/naiba/gopappy"
	"git.cm/naiba/gopappy/util/aliyun"
	"git.cm/naiba/gopappy/util/cn4"
	"git.cm/naiba/gopappy/util/ename"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
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
	en, _ := ename.Domains(o)
	cn, _ := cn4.Domains(o)
	ali, _ := aliyun.Domains(o)
	en = append(en, cn...)
	en = append(en, ali...)
	convertPrices(en)
	sort.Sort(gopappy.SortDomain(en))
	c.JSON(http.StatusOK, en)
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
	type Params struct {
		Platforms map[int]string `json:"platforms"`
		TLDs      map[int]string `json:"tlds"`
		Tags      map[int]string `json:"tags"`
	}
	c.JSON(http.StatusOK, Params{
		Platforms: gopappy.Platform,
		TLDs:      gopappy.TLDs,
		Tags:      gopappy.Tags,
	})
}

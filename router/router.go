/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gopappy"
	"git.cm/naiba/gopappy/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"git.cm/naiba/gopappy/util/aliyun"
)

func init() {
	service.Init()
}

func RunWeb() {
	g := gin.Default()
	g.Use(func(c *gin.Context) {
		abort := func() {
			c.String(http.StatusForbidden, "被安全策略拦截")
			c.Abort()
		}
		if len(c.Request.UserAgent()) < 20 {
			abort()
		}
		if len(c.Request.Header["Accept-Language"]) == 0 || len(c.Request.Header["Accept-Language"][0]) < 2 {
			abort()
		}
		if len(c.Request.Header["Referer"]) == 0 || len(c.Request.Header["Referer"][0]) < 5 {
			abort()
		}
	})
	g.GET("/test", func(c *gin.Context) {
		log.Println(aliyun.Domains(gopappy.Option{
			Page:    1,
			Keyword: "gh",
			KwPos:   1,
			TLDs: []int{
				1,
			},
			Tag:       6,
			MaxLength: 3,
			Order:     1,
			Sort:      1,
		}))
	})
	g.Run("0.0.0.0:8000")
}

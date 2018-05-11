/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gopappy"
	"git.cm/naiba/gopappy/service"
	"git.cm/naiba/gopappy/util/cn4"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	service.Init()
}

func RunWeb() {
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		log.Println(cn4.Domains(gopappy.Option{
			Page: 1,
			TLDs: []int{
				1,
				2,
			},
			Keyword: "livelive",
			Order:   1,
			Sort:    2,
		}))
	})
	g.Run("0.0.0.0:8000")
}

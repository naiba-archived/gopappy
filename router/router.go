/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package router

import (
	"git.cm/naiba/gopappy/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
}

func RunWeb() {
	g := gin.Default()
	g.Use(middleware.Cross)
	g.Use(middleware.Safe)
	serveAPI(g)
	g.Run("127.0.0.1:3010")
}

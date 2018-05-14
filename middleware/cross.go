/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cross(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Method", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	if c.Request.Method == http.MethodOptions {
		c.Status(http.StatusOK)
		c.Abort()
	}
}

/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Safe(c *gin.Context) {
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
}

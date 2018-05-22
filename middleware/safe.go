/*
 * Copyright (c) 2018, 奶爸<1@5.nu>
 * All rights reserved.
 */

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Safe(c *gin.Context) {
	abort := func() {
		c.String(http.StatusForbidden, "被安全策略拦截")
		c.Abort()
	}
	if len(c.Request.UserAgent()) < 20 {
		abort()
		return
	}
	if len(c.Request.Header["Accept-Language"]) == 0 || len(c.Request.Header["Accept-Language"][0]) < 2 {
		abort()
		return
	}
	if len(c.Request.Header["Referer"]) == 0 || len(c.Request.Header["Referer"][0]) < 5 {
		abort()
		return
	}
}

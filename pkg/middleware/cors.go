package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		if origin != "" {
			// 允许的源
			c.Header("Access-Control-Allow-Origin", origin)
			// 允许的请求方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			// 允许的请求头
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, X-CSRF-Token, Token, session, X-Requested-With")
			// 允许暴露的响应头
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// 是否允许携带凭证（cookies等）
			c.Header("Access-Control-Allow-Credentials", "true")
			// 预检请求的有效期（秒）
			c.Header("Access-Control-Max-Age", "86400")
		}
		
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// CORSWithConfig 带配置的跨域中间件
func CORSWithConfig(allowOrigins []string, allowMethods []string, allowHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		// 检查origin是否在允许列表中
		allowed := false
		for _, allowOrigin := range allowOrigins {
			if allowOrigin == "*" || allowOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			
			// 设置允许的方法
			methods := "GET, POST, PUT, DELETE, OPTIONS"
			if len(allowMethods) > 0 {
				methods = ""
				for i, m := range allowMethods {
					if i > 0 {
						methods += ", "
					}
					methods += m
				}
			}
			c.Header("Access-Control-Allow-Methods", methods)
			
			// 设置允许的头
			headers := "Authorization, Content-Length, Content-Type, X-Requested-With"
			if len(allowHeaders) > 0 {
				headers = ""
				for i, h := range allowHeaders {
					if i > 0 {
						headers += ", "
					}
					headers += h
				}
			}
			c.Header("Access-Control-Allow-Headers", headers)
			
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}
		
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}


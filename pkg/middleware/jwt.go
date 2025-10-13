package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	pkgRedis "goweb/pkg/redis"
	"goweb/pkg/response"
)

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return JWTAuthWithRedis(secret, nil)
}

// JWTAuthWithRedis JWT认证中间件（支持token黑名单）
func JWTAuthWithRedis(secret string, redisClient *pkgRedis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 检查token是否在黑名单中
		if redisClient != nil {
			ctx := context.Background()
			blacklistKey := fmt.Sprintf("token:blacklist:%s", tokenString)
			exists, err := redisClient.Exists(ctx, blacklistKey)
			if err == nil && exists {
				response.Unauthorized(c, "认证令牌已失效，请重新登录")
				c.Abort()
				return
			}
		}

		// 解析JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			response.Unauthorized(c, "认证令牌解析失败: "+err.Error())
			c.Abort()
			return
		}

		if !token.Valid {
			response.Unauthorized(c, "认证令牌无效")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "认证令牌解析失败")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		c.Next()
	}
}

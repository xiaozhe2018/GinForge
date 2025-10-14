package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WebSocketAuth WebSocket 认证中间件
// 从 URL 参数中获取 token 进行验证
func WebSocketAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 URL 参数获取 token（WebSocket 无法使用 Header）
		token := c.Query("token")
		if token == "" {
			// 尝试从 Header 获取（用于 HTTP 升级前的验证）
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}
		
		if token == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}
		
		// 解析 JWT Token
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		
		if err != nil || !parsedToken.Valid {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}
		
		// 提取用户信息
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		
		c.Next()
	}
}


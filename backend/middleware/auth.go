package middleware

import (
	"strings"

	"qiniudemo/config"
	"qiniudemo/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			model.Error(c, 401, "请先登录")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			model.Error(c, 401, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			model.Error(c, 401, "令牌无效")
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			model.Error(c, 401, "令牌无效")
			c.Abort()
			return
		}

		c.Set("user_id", int64(userID))
		c.Next()
	}
}

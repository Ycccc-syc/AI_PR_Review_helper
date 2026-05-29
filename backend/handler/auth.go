package handler

import (
	"database/sql"
	"time"

	"qiniudemo/config"
	"qiniudemo/model"
	"qiniudemo/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Register(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供有效的用户名（2-32位）和密码（6-128位）")
			return
		}

		user, err := repository.CreateUser(db, req.Username, req.Password)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}

		model.Success(c, model.AuthResponse{
			Token: "",
			User:  *user,
		})
	}
}

func Login(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供用户名和密码")
			return
		}

		user, err := repository.GetUserByUsername(db, req.Username)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}

		if !repository.VerifyPassword(user.Password, req.Password) {
			model.Error(c, -1, "用户名或密码错误")
			return
		}

		token, err := generateToken(user.ID, cfg.JWTSecret)
		if err != nil {
			model.Error(c, -1, "生成令牌失败")
			return
		}

		model.Success(c, model.AuthResponse{
			Token: token,
			User:  *user,
		})
	}
}

func generateToken(userID int64, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

package repository

import (
	"database/sql"
	"fmt"

	"qiniudemo/model"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB, username, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}
	var user model.User
	err = db.QueryRow(
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at",
		username, string(hash),
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return nil, fmt.Errorf("用户名已被注册")
		}
		return nil, fmt.Errorf("注册失败: %w", err)
	}
	return &user, nil
}

func GetUserByUsername(db *sql.DB, username string) (*model.User, error) {
	var user model.User
	err := db.QueryRow(
		"SELECT id, username, password, created_at FROM users WHERE username=$1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

func VerifyPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func isDuplicate(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate") || contains(err.Error(), "唯一") || contains(err.Error(), "unique"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

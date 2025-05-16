package opsuser

import (
	"errors"
	"fmt"
	"os"
	"os/user"
)

// UserInfo содержит информацию о пользователе
type UserInfo struct {
	UID      string
	HomeDir  string
	Username string
}

// GetByUsername возвращает информацию о пользователе по его имени
func GetByUsername(username string) (*UserInfo, error) {
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("lookup user %s: %w", username, err)
	}

	return &UserInfo{
		UID:      u.Uid,
		HomeDir:  u.HomeDir,
		Username: u.Username,
	}, nil
}

// CheckExists проверяет, существует ли пользователь по его имени
func CheckExists(username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username cannot be empty")
	}

	_, err := user.Lookup(username)
	if err != nil {
		var unknownUserError user.UnknownUserError
		if errors.As(err, &unknownUserError) {
			return false, err
		}
		return false, fmt.Errorf("check user %s: %w", username, err)
	}

	return true, nil
}

// GetCurrent возвращает информацию о текущем пользователе
func GetCurrent() (*UserInfo, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}

	return &UserInfo{
		UID:      u.Uid,
		HomeDir:  u.HomeDir,
		Username: u.Username,
	}, nil
}

// IsRoot проверяет, является ли текущий пользователь суперпользователем (root)
func IsRoot() bool {
	return os.Geteuid() == 0
}

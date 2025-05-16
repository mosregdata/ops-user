package opsuser

import (
	"errors"
	"os/user"
	"reflect"
	"testing"
)

// TestGetByUsername проверяет получение информации о пользователе по имени
func TestGetByUsername(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("failed to get current user: %v", err)
	}

	info, err := GetByUsername(currentUser.Username)
	if err != nil {
		t.Errorf("GetByUsername(%s) failed: %v", currentUser.Username, err)
	}

	if info == nil {
		t.Fatal("GetByUsername returned nil info")
	}

	if info.Username != currentUser.Username {
		t.Errorf("expected Username %s, got %s", currentUser.Username, info.Username)
	}
	if info.UID != currentUser.Uid {
		t.Errorf("expected UID %s, got %s", currentUser.Uid, info.UID)
	}
	if info.HomeDir != currentUser.HomeDir {
		t.Errorf("expected HomeDir %s, got %s", currentUser.HomeDir, info.HomeDir)
	}

	_, err = GetByUsername("")
	if err == nil {
		t.Error("GetByUsername with empty username should return error")
	}

	_, err = GetByUsername("nonexistentuser123")
	if err == nil {
		t.Error("GetByUsername with nonexistent user should return error")
	}
}

// TestGetCurrent проверяет получение информации о текущем пользователе
func TestGetCurrent(t *testing.T) {
	info, err := GetCurrent()
	if err != nil {
		t.Errorf("GetCurrent failed: %v", err)
	}

	if info == nil {
		t.Fatal("GetCurrent returned nil info")
	}

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("failed to get current user: %v", err)
	}

	if !reflect.DeepEqual(info, &UserInfo{
		UID:      currentUser.Uid,
		HomeDir:  currentUser.HomeDir,
		Username: currentUser.Username,
	}) {
		t.Errorf("GetCurrent returned unexpected info: got %+v", info)
	}
}

// Mock для user.Lookup, чтобы избежать реальных системных вызовов
type mockUserLookup func(string) (*user.User, error)

func (m mockUserLookup) Lookup(username string) (*user.User, error) {
	return m(username)
}

func TestCheckExists(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		lookupMock func(string) (*user.User, error)
		wantExists bool
		wantErr    bool
	}{
		{
			name:     "Empty username",
			username: "",
			lookupMock: func(_ string) (*user.User, error) {
				return nil, errors.New("should not be called")
			},
			wantExists: false,
			wantErr:    true,
		},
		{
			name:     "Existing user",
			username: "root",
			lookupMock: func(_ string) (*user.User, error) {
				return &user.User{Username: "root"}, nil
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:     "Non-existing user",
			username: "bob",
			lookupMock: func(_ string) (*user.User, error) {
				return nil, user.UnknownUserError("bob")
			},
			wantExists: false,
			wantErr:    true,
		},
		{
			name:     "Lookup error",
			username: "charlie",
			lookupMock: func(_ string) (*user.User, error) {
				return nil, errors.New("lookup failed")
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalLookup := userLookup
			userLookup = mockUserLookup(tt.lookupMock)
			defer func() { userLookup = originalLookup }()

			gotExists, err := CheckExists(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotExists != tt.wantExists {
				t.Errorf("CheckExists() exists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

var userLookup = user.Lookup

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_AddPushToken(t *testing.T) {
	t.Run("add push token to user", func(t *testing.T) {
		user := &User{
			ID:        "id",
			PushToken: map[string]bool{},
		}
		newPushToken := "new"
		user.AddPushToken(newPushToken)
		assert.Equal(t, 1, len(user.PushToken))
		assert.Equal(t, true, user.PushToken[newPushToken])
	})
	t.Run("nil push token map", func(t *testing.T) {
		user := &User{
			ID:        "id",
			PushToken: nil,
		}
		newPushToken := "new"
		user.AddPushToken(newPushToken)
		assert.Equal(t, 1, len(user.PushToken))
		assert.Equal(t, true, user.PushToken[newPushToken])
	})
	t.Run("add duplicate push token", func(t *testing.T) {
		duplicatedPushToken := "duplicate"
		user := &User{
			ID:        "id",
			PushToken: map[string]bool{duplicatedPushToken: true},
		}
		user.AddPushToken(duplicatedPushToken)
		assert.Equal(t, 1, len(user.PushToken))
		assert.Equal(t, true, user.PushToken[duplicatedPushToken])
	})
}

func TestUser_RemovePushToken(t *testing.T) {
	t.Run("remove push token", func(t *testing.T) {
		removeToken := "remove"
		user := &User{
			ID:        "id",
			PushToken: map[string]bool{removeToken: true},
		}
		user.RemovePushToken(removeToken)
		assert.Equal(t, 0, len(user.PushToken))
	})
	t.Run("nil push token map", func(t *testing.T) {
		removeToken := "token"
		user := &User{
			ID:        "id",
			PushToken: nil,
		}
		user.RemovePushToken(removeToken)
		assert.Equal(t, 0, len(user.PushToken))
	})
	t.Run("remove doen't existed token", func(t *testing.T) {
		notExistedToken := "not exist"
		user := &User{
			ID:        "id",
			PushToken: map[string]bool{},
		}
		user.RemovePushToken(notExistedToken)
		assert.Equal(t, 0, len(user.PushToken))
	})
}

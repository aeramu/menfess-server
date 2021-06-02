package service

import (
	"context"
	"errors"
	"testing"

	"github.com/aeramu/menfess-server/internal/notification/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initTest() (
	service Service,
	repoMock *MockRepository,
	pushMock *MockPushServiceClient,
) {
	repoMock = new(MockRepository)
	pushMock = new(MockPushServiceClient)
	service = NewService(repoMock, pushMock)

	return
}

func Test_service_AddPushToken(t *testing.T) {
	ctx := context.Background()
	req := AddPushTokenReq{
		ID:        "id",
		PushToken: "new token",
	}
	user := User{
		ID:        "id",
		PushToken: map[string]bool{"token1": true, "token2": true},
	}
	t.Run("error find by id", func(t *testing.T) {
		service, repoMock, _ := initTest()
		expectedErr := errors.New("repo error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(nil, expectedErr)
		err := service.AddPushToken(ctx, req)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("user not found", func(t *testing.T) {
		service, repoMock, _ := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(nil, nil)
		err := service.AddPushToken(ctx, req)
		assert.Equal(t, constants.ErrUserNotFound, err)
	})
	t.Run("error repo save", func(t *testing.T) {
		service, repoMock, _ := initTest()
		expectedErr := errors.New("repo error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(&user, nil)
		repoMock.On("Save",
			mock.Anything,
			mock.MatchedBy(func(u User) bool {
				assert.Equal(t, 3, len(u.PushToken))
				return true
			})).
			Return(expectedErr)
		err := service.AddPushToken(ctx, req)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("success", func(t *testing.T) {
		service, repoMock, _ := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(&user, nil)
		repoMock.On("Save",
			mock.Anything,
			mock.MatchedBy(func(u User) bool {
				assert.Equal(t, 3, len(u.PushToken))
				return true
			})).
			Return(nil)
		err := service.AddPushToken(ctx, req)
		assert.NoError(t, err)
	})
}

func Test_service_RemovePushToken(t *testing.T) {
	ctx := context.Background()
	req := RemovePushTokenReq{
		ID:        "id",
		PushToken: "token1",
	}
	user := User{
		ID:        "id",
		PushToken: map[string]bool{"token1": true, "token2": true},
	}
	t.Run("error find by id", func(t *testing.T) {
		service, repoMock, _ := initTest()
		expectedErr := errors.New("repo error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(nil, expectedErr)
		err := service.RemovePushToken(ctx, req)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("user not found", func(t *testing.T) {
		service, repoMock, _ := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(nil, nil)
		err := service.RemovePushToken(ctx, req)
		assert.Equal(t, constants.ErrUserNotFound, err)
	})
	t.Run("error repo save", func(t *testing.T) {
		service, repoMock, _ := initTest()
		expectedErr := errors.New("repo error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(&user, nil)
		repoMock.On("Save",
			mock.Anything,
			mock.MatchedBy(func(u User) bool {
				assert.Equal(t, 1, len(u.PushToken))
				return true
			})).
			Return(expectedErr)
		err := service.RemovePushToken(ctx, req)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("success", func(t *testing.T) {
		service, repoMock, _ := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.ID, id)
				return true
			})).
			Return(&user, nil)
		repoMock.On("Save",
			mock.Anything,
			mock.MatchedBy(func(u User) bool {
				assert.Equal(t, 1, len(u.PushToken))
				return true
			})).
			Return(nil)
		err := service.RemovePushToken(ctx, req)
		assert.NoError(t, err)
	})
}

func Test_service_SendNotification(t *testing.T) {
	ctx := context.Background()
	req := SendNotificationReq{
		Title:  "title",
		Body:   "body",
		UserID: "user1",
		Data:   "json here",
	}
	user := User{
		ID:        "user1",
		PushToken: map[string]bool{"token1": true, "token2": true},
	}
	t.Run("error find by id", func(t *testing.T) {
		service, repoMock, _ := initTest()
		expectedErr := errors.New("repo error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.Anything).
			Return(nil, expectedErr)
		err := service.SendNotification(ctx, req)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("user not found", func(t *testing.T) {
		service, repoMock, _ := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.UserID, id)
				return true
			})).
			Return(nil, nil)
		err := service.SendNotification(ctx, req)
		assert.Equal(t, constants.ErrUserNotFound, err)
	})
	t.Run("error send to push service", func(t *testing.T) {
		service, repoMock, pushMock := initTest()
		expectedErr := errors.New("something error")
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(id string) bool {
				assert.Equal(t, req.UserID, id)
				return true
			})).
			Return(&user, nil)
		pushMock.On("Send",
			mock.Anything,
			mock.AnythingOfType("map[string]bool"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(expectedErr)
		err := service.SendNotification(ctx, req)
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		service, repoMock, pushMock := initTest()
		repoMock.On("FindByID",
			mock.Anything,
			mock.MatchedBy(func(s string) bool {
				assert.Equal(t, req.UserID, s)
				return true
			})).
			Return(&user, nil)
		pushMock.On("Send",
			mock.Anything,
			mock.MatchedBy(func(m map[string]bool) bool {
				assert.Equal(t, user.PushToken, m)
				return true
			}),
			mock.MatchedBy(func(s string) bool {
				assert.Equal(t, req.Title, s)
				return true
			}),
			mock.MatchedBy(func(s string) bool {
				assert.Equal(t, req.Body, s)
				return true
			}),
			mock.MatchedBy(func(s string) bool {
				assert.Equal(t, req.Data, s)
				return true
			})).
			Return(nil)
		err := service.SendNotification(ctx, req)
		assert.NoError(t, err)
	})
}

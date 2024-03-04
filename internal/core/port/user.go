package port

import (
	"context"
	"mqtt_2_sms_webserver/internal/core/entity"
)

type UserRepository interface {
	GetOne(ctx context.Context, user *entity.User) (entity.User, error)
}

type UserService interface {
	GetAllUsers() ([]*entity.User, error)
	GetUser(ctx context.Context, user *entity.User) (entity.User, error)
}

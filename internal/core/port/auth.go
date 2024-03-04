package port

import (
	"context"
	"mqtt_2_sms_webserver/internal/core/entity"
)

type AuthService interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CheckPassword(ctx context.Context, user *entity.User, password string) error
}

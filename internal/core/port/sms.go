package port

import "context"

type SMSRepository interface {
}

type SMSService interface {
	SendSMS(ctx context.Context, phone_number, message string) error
}

package request

type SMSSendRequest struct {
	Phone   string `json:"phone" binding:"required,e164"`
	Message string `json:"message" binding:"required"`
}

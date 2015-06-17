package git

type HashAndMessage struct {
	Hash string `json:"hash" binding:"required"`
	Message string `json:"message" binding:"required"`
}
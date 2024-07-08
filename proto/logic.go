package proto

type LoginRequest struct {
	Name     string
	Password string
}

type LoginResponse struct {
	Code      int
	AuthToken string
}

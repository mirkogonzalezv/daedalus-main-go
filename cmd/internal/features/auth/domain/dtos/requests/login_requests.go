package requests

type LoginRequest struct {
	Email     string
	Password  string
	IP        string
	UserAgent string
}

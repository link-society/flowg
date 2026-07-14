package schemas

// LoginRequest carries the credentials presented for authentication.
type LoginRequest struct {
	// Username is the account name to authenticate.
	Username string `json:"username" required:"true"`
	// Password is the account's password.
	Password string `json:"password" required:"true"`
}

// LoginResponse carries the session token issued on success.
type LoginResponse struct {
	// Success reports whether authentication succeeded.
	Success bool `json:"success"`
	// Token is a JWT proving the caller's identity on subsequent requests.
	Token string `json:"token"`
}

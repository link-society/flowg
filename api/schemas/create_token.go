package schemas

// CreateTokenRequest is empty: the token is issued for the calling user.
type CreateTokenRequest struct{}

// CreateTokenResponse carries the newly issued personal access token.
type CreateTokenResponse struct {
	// Success reports whether the token was created.
	Success bool `json:"success"`
	// Token is the secret value, returned only once at creation time.
	Token string `json:"token"`
	// TokenUUID identifies the token for later listing or deletion.
	TokenUUID string `json:"token_uuid"`
}

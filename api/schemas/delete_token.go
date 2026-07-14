package schemas

// DeleteTokenRequest identifies the token to revoke.
type DeleteTokenRequest struct {
	// TokenUUID is the identifier of the caller's token to delete.
	TokenUUID string `path:"token-uuid" format:"uuid"`
}

// DeleteTokenResponse reports the outcome of the revocation.
type DeleteTokenResponse struct {
	// Success reports whether the token was revoked.
	Success bool `json:"success"`
}

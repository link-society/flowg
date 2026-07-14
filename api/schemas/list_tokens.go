package schemas

// ListTokensRequest is empty: tokens are listed for the calling user.
type ListTokensRequest struct{}

// ListTokensResponse carries the identifiers of the caller's tokens.
type ListTokensResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// TokenUUIDs identifies each of the caller's tokens; the secret values are
	// never returned.
	TokenUUIDs []string `json:"token_uuids"`
}

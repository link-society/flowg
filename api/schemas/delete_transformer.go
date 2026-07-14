package schemas

// DeleteTransformerRequest identifies the transformer to remove.
type DeleteTransformerRequest struct {
	// Transformer is the name of the transformer to delete.
	Transformer string `path:"transformer" minLength:"1"`
}

// DeleteTransformerResponse reports the outcome of the deletion.
type DeleteTransformerResponse struct {
	// Success reports whether the transformer was removed.
	Success bool `json:"success"`
}

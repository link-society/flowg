package schemas

// GetTransformerRequest identifies the transformer to retrieve.
type GetTransformerRequest struct {
	// Transformer is the name of the transformer to read.
	Transformer string `path:"transformer" minLength:"1"`
}

// GetTransformerResponse carries the source of the requested transformer.
type GetTransformerResponse struct {
	// Success reports whether the transformer was found and returned.
	Success bool `json:"success"`
	// Script is the VRL source code of the transformer.
	Script string `json:"script"`
}

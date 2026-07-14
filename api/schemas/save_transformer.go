package schemas

// SaveTransformerRequest carries the transformer name and its new source.
type SaveTransformerRequest struct {
	// Transformer is the name of the transformer to create or overwrite.
	Transformer string `path:"transformer" minLength:"1"`
	// Script is the VRL source code to store under that name.
	Script string `json:"script" required:"true"`
}

// SaveTransformerResponse reports the outcome of the save.
type SaveTransformerResponse struct {
	// Success reports whether the transformer was persisted.
	Success bool `json:"success"`
}

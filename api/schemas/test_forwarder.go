package schemas

// TestForwarderRequest carries the forwarder name and a sample record to send
// through it.
type TestForwarderRequest struct {
	// Forwarder is the name of the stored forwarder to exercise.
	Forwarder string `path:"forwarder" minLength:"1"`
	// Record is the sample log record to forward.
	Record map[string]string `json:"record" required:"true"`
}

// TestForwarderResponse reports whether the trial delivery succeeded.
type TestForwarderResponse struct {
	// Success reports whether the record was delivered.
	Success bool `json:"success"`
}

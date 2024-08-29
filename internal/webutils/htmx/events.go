package htmx

type ModalOpenEvent struct{}

type ModalCloseEvent struct {
	After string `json:"after"`
}

type ToastEvent struct {
	Messages []string `json:"messages"`
}

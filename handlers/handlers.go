package handlers

type R struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data string `json:"data"`
}

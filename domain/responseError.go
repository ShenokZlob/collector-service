package domain

type ResponseErr struct {
	Message string `json:"message"`
	Status  int    `json:"status,omitempty"`
}

func (r *ResponseErr) Error() string {
	if r == nil {
		return "ResponseErr equals nil!"
	}
	return r.Message
}

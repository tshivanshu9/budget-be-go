package custom_errors

type NotFoundError struct {
	Message string
}

func (r *NotFoundError) Error() string {
	return r.Message
}

func NewNotFoundError(message string) *NotFoundError {
	if message == "" {
		message = "Requested resource not found"
	}
	return &NotFoundError{
		Message: message,
	}
}

package apperror

type AppError struct {
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NotFound(resource string) error {
	return AppError{
		Message: resource + " not found",
	}
}

func Failed(resource string) error {
	return AppError{
		Message: "failed to " + resource,
	}
}

func Internal() error {
	return AppError{
		Message: "Internal server error",
	}
}
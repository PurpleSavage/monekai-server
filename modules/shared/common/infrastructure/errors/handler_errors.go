package globalerrors

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	appErr, ok := err.(*AppError)
	return ok && appErr.Status == 404
}
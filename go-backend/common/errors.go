package common

func WrapError(message string, e error) ErrorWrapper {
	return ErrorWrapper{Message: message, InnerError: e}
}

type ErrorWrapper struct {
	Message    string
	InnerError error
}

var _ error = ErrorWrapper{}

func (e ErrorWrapper) Error() string {
	return e.Message + " [CAUSED BY]\n" + e.InnerError.Error()
}

package townsta

type HTTPError struct {
	Err     error
	Message string
	Code    int
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

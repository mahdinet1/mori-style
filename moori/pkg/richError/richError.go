package richError

type Kind int

const (
	InvalidCode Kind = iota + 1
	ForbiddenCode
	NotFoundCode
	AlreadyExistsCode
	UnexpectedCode
)

type RichError struct {
	message      string
	operation    string
	meta         map[string]interface{}
	wrappedError error
	code         Kind
}

func New(op string) RichError {
	return RichError{operation: op}
}

func (e RichError) RetrieveAncestorCode() Kind {
	if e.code != 0 {
		return e.code
	}
	re, ok := e.wrappedError.(RichError)
	if !ok {
		return 0
	}
	return re.RetrieveAncestorCode()

}

func (e RichError) RetrieveAncestorMsg() string {
	if e.code != 0 {
		return e.Error()
	}
	re, ok := e.wrappedError.(RichError)
	if !ok {
		return re.Error()
	}
	return re.RetrieveAncestorMsg()

}
func (e RichError) RetrieveCode() Kind {
	return e.code
}

func (e RichError) RetrieveMsg() string {
	return e.message
}
func (e RichError) Error() string {
	return e.message
}

func (e RichError) SetMessage(msg string) RichError {
	e.message = msg
	return e
}
func (e RichError) SetMeta(meta map[string]interface{}) RichError {
	e.meta = meta
	return e
}
func (e RichError) SetWrappedError(err error) RichError {
	e.wrappedError = err
	return e
}
func (e RichError) SetCode(code Kind) RichError {
	e.code = code
	return e
}
func (e RichError) RetrieveOperation() string {
	return e.operation
}

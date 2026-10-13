package apperror

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindUnauthorized
	KindForbidden
	KindNotFound
	KindConflict
	KindInternal
)

type Error struct {
	Kind    Kind
	Message string // safe to show to the client
	Details string // optional extra info, safe to show
	Err     error  // underlying cause: logged, never sent to the client
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

// Is lets errors.Is match sentinel errors even after WithDetails/Wrap copies them.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Kind == e.Kind && t.Message == e.Message
}

func (e *Error) WithDetails(details string) *Error {
	c := *e
	c.Details = details
	return &c
}

func (e *Error) Wrap(err error) *Error {
	c := *e
	c.Err = err
	return &c
}

func New(kind Kind, message string) *Error { return &Error{Kind: kind, Message: message} }

func Invalid(message string) *Error      { return New(KindInvalid, message) }
func Unauthorized(message string) *Error { return New(KindUnauthorized, message) }
func Forbidden(message string) *Error    { return New(KindForbidden, message) }
func NotFound(message string) *Error     { return New(KindNotFound, message) }
func Conflict(message string) *Error     { return New(KindConflict, message) }
func Internal(message string) *Error     { return New(KindInternal, message) }

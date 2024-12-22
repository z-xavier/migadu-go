package migadu

type MigaduError struct {
	httpStatusCode int
	error          string
}

func (e *MigaduError) Error() string {
	return e.error
}

func NewMigaduError(error string) *MigaduError {
	return &MigaduError{}
}

const (
	ErrorUnauthorized = "the token is invalid or expired"
)

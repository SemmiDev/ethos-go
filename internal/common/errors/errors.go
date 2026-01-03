package errors

type SlugError struct {
	error    string
	slug     string
	httpCode int
}

func (s SlugError) Error() string       { return s.error }
func (s SlugError) Slug() string        { return s.slug }
func (s SlugError) HTTPStatusCode() int { return s.httpCode }

func NewIncorrectInputError(msg, slug string) SlugError {
	return SlugError{error: msg, slug: slug, httpCode: 400}
}

func NewNotFoundError(msg, slug string) SlugError {
	return SlugError{error: msg, slug: slug, httpCode: 404}
}

func NewUnauthorizedError(msg, slug string) SlugError {
	return SlugError{error: msg, slug: slug, httpCode: 403}
}

func NewConflictError(msg, slug string) SlugError {
	return SlugError{error: msg, slug: slug, httpCode: 409}
}

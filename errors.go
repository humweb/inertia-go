package inertia

import "errors"

var (
	// ErrInvalidContextProps error.
	ErrInvalidContextProps = errors.New("inertia: could not convert context props to map")

	// ErrInvalidContextViewData error.
	ErrInvalidContextViewData = errors.New("inertia: could not convert context view data to map")

	// ErrBadSsrStatusCode error.
	ErrBadSsrStatusCode = errors.New("inertia: bad ssr status code >= 400")

	// ErrBadSsrStatusCode error.
	ErrRawTemplateFunc = errors.New("inertia: error with raw template func")
)

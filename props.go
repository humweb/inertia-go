package inertia

import (
	"fmt"
	"net/http"
	"strings"
)

// LazyProp is a property value that will only evaluate when needed.
//
// https://inertiajs.com/partial-reloads
type LazyProp func() (any, error)

func (i *Inertia) PrepareProps(r *http.Request, component string, props Props) (Props, error) {

	isPartial := r.Header.Get(Headers.PartialComponent) == component

	// Merge props and shared props
	for k, v := range i.SharedProps {
		if _, ok := props[k]; !ok {
			props[k] = v
		}
	}

	// Add props from context to the result.
	contextProps := r.Context().Value(ContextKeyProps)

	if contextProps != nil {
		contextProps, ok := contextProps.(Props)
		if !ok {
			return nil, ErrInvalidContextProps
		}

		for key, value := range contextProps {
			props[key] = value
		}
	}

	// Get props keys to return. If len == 0, then return all.
	partials := r.Header.Get(Headers.PartialOnly)
	data := strings.Split(partials, ",")
	only := make(map[string]struct{}, len(data))

	if partials != "" && isPartial {
		for _, value := range data {
			only[value] = struct{}{}
		}
	}

	// Filter props.
	if len(only) > 0 {
		// While making partials requests:
		// Use the `only` property to specify which data the server should return.
		for key := range props {
			if _, ok := only[key]; !ok {
				delete(props, key)
			}
		}
	} else {
		// Lazy props should only be evaluated when required using the `only` property
		for key, val := range props {
			if _, ok := val.(LazyProp); ok {
				delete(props, key)
			}
		}
	}

	// Resolve props values.
	for key, val := range props {
		val, err := ResolvePropVal(val)
		if err != nil {
			return nil, fmt.Errorf("resolve prop value: %w", err)
		}
		props[key] = val
	}

	return props, nil
}

func ResolvePropVal(val any) (any, error) {
	var err error

	if closure, ok := val.(func() (any, error)); ok {

		if val, err = closure(); err != nil {

			return nil, fmt.Errorf("closure prop resolving: %w", err)

		}
	} else if lazy, ok := val.(LazyProp); ok {
		val, err = lazy()

		if err != nil {
			return nil, fmt.Errorf("lazy prop resolving: %w", err)
		}
	}

	return val, nil
}

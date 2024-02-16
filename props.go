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

func (i *Inertia) prepareProps(r *http.Request, component string, props Props) (Props, error) {
	result := make(Props)

	// Add shared props to the result.
	for key, val := range i.sharedProps {
		result[key] = val
	}

	// Add props from context to the result.
	contextProps := r.Context().Value(ContextKeyProps)

	if contextProps != nil {
		contextProps, ok := contextProps.(Props)
		if !ok {
			return nil, ErrInvalidContextProps
		}

		for key, value := range contextProps {
			result[key] = value
		}
	}

	// Add passed props to the result.
	for key, val := range props {
		result[key] = val
	}

	// Get props keys to return. If len == 0, then return all.
	partial := r.Header.Get("X-Inertia-Partial-Data")
	data := strings.Split(partial, ",")
	only := make(map[string]struct{}, len(data))

	if partial != "" && r.Header.Get("X-Inertia-Partial-Component") == component {
		for _, value := range data {
			only[value] = struct{}{}
		}
	}

	// Filter props.
	if len(only) > 0 {
		// While making partial requests:
		// Use the `only` property to specify which data the server should return.
		for key := range result {
			if _, ok := only[key]; !ok {
				delete(result, key)
			}
		}
	} else {
		// Lazy props should only be evaluated when required using the `only` property
		for key, val := range result {
			if _, ok := val.(LazyProp); ok {
				delete(result, key)
			}
		}
	}

	// Resolve props values.
	for key, val := range result {
		val, err := resolvePropVal(val)
		if err != nil {
			return nil, fmt.Errorf("resolve prop value: %w", err)
		}
		result[key] = val
	}

	return result, nil
}

func resolvePropVal(val any) (_ any, err error) {
	if closure, ok := val.(func() (any, error)); ok {
		val, err = closure()
		if err != nil {
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

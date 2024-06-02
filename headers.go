package inertia

type HeaderTypes struct {
	Inertia          string
	ErrorBag         string
	Location         string
	Version          string
	PartialComponent string
	PartialOnly      string
	PartialExcept    string
}

var Headers = HeaderTypes{
	Inertia:          "X-Inertia",
	ErrorBag:         "X-Inertia-Error-Bag",
	Location:         "X-Inertia-Location",
	Version:          "X-Inertia-Version",
	PartialComponent: "X-Inertia-Partial-Component",
	PartialOnly:      "X-Inertia-Partial-Data",
	PartialExcept:    "X-Inertia-Partial-Except",
}

package inertia

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

// Inertia type.
type Inertia struct {
	Url           string
	rootTemplate  string
	version       string
	SharedProps   Props
	SharedFuncMap template.FuncMap
	templateFS    fs.FS
	SsrURL        string
	SsrClient     *http.Client
}

// New function.
func New(url, rootTemplate, version string) *Inertia {
	return &Inertia{
		Url:           url,
		rootTemplate:  rootTemplate,
		version:       version,
		SharedFuncMap: template.FuncMap{"marshal": Marshal, "raw": Raw},
		SharedProps:   Props{},
	}
}

// NewWithFS function.
func NewWithFS(url, rootTemplate, version string, templateFS fs.FS) *Inertia {
	i := New(url, rootTemplate, version)
	i.templateFS = templateFS

	return i
}

// IsSsrEnabled function.
func (i *Inertia) IsSsrEnabled() bool {
	return i.SsrURL != "" && i.SsrClient != nil
}

// EnableSsr function.
func (i *Inertia) EnableSsr(ssrURL string) {
	i.SsrURL = ssrURL
	i.SsrClient = &http.Client{}
}

// EnableSsrWithDefault function.
func (i *Inertia) EnableSsrWithDefault() {
	i.EnableSsr("http://127.0.0.1:13714")
}

// DisableSsr function.
func (i *Inertia) DisableSsr() {
	i.SsrURL = ""
	i.SsrClient = nil
}

// Share function.
func (i *Inertia) Share(key string, value any) {
	i.SharedProps[key] = value
}

// ShareFunc function.
func (i *Inertia) ShareFunc(key string, value any) {
	i.SharedFuncMap[key] = value
}

// WithProp function.
func (i *Inertia) WithProp(ctx context.Context, key string, value any) context.Context {
	contextProps := ctx.Value(ContextKeyProps)

	if contextProps != nil {
		contextProps, ok := contextProps.(Props)
		if ok {
			contextProps[key] = value

			return context.WithValue(ctx, ContextKeyProps, contextProps)
		}
	}

	return context.WithValue(ctx, ContextKeyProps, Props{
		key: value,
	})
}

// WithProps appends props values to the passed context.Context.
func (i *Inertia) WithProps(ctx context.Context, props Props) context.Context {
	if ctxData := ctx.Value(ContextKeyProps); ctxData != nil {
		ctxData, ok := ctxData.(Props)

		if ok {
			for key, val := range props {
				ctxData[key] = val
			}

			return context.WithValue(ctx, ContextKeyProps, ctxData)
		}
	}

	return context.WithValue(ctx, ContextKeyProps, props)
}

// WithViewData function.
func (i *Inertia) WithViewData(ctx context.Context, key string, value any) context.Context {
	contextViewData := ctx.Value(ContextKeyViewData)

	if contextViewData != nil {
		contextViewData, ok := contextViewData.(Props)
		if ok {
			contextViewData[key] = value

			return context.WithValue(ctx, ContextKeyViewData, contextViewData)
		}
	}

	return context.WithValue(ctx, ContextKeyViewData, Props{
		key: value,
	})
}

// Render function.
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props Props) error {

	preparedProps, err := i.PrepareProps(r, component, props)

	page := &Page{
		Component: component,
		URL:       r.RequestURI,
		Version:   i.version,
		Props:     preparedProps,
	}

	// Inertia request
	if i.isInertiaRequest(r) {
		js, err := json.Marshal(page)
		if err != nil {
			return err
		}

		w.Header().Set("Vary", "Accept")
		w.Header().Set("X-Inertia", "true")
		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(js)
		if err != nil {
			return err
		}

		return nil
	}

	// View data
	contextViewData := r.Context().Value(ContextKeyViewData)
	viewData := make(Props)

	if contextViewData != nil {
		contextViewData, ok := contextViewData.(Props)
		if !ok {
			return ErrInvalidContextViewData
		}

		for key, value := range contextViewData {
			viewData[key] = value
		}
	}
	viewData["page"] = page

	if i.IsSsrEnabled() {
		ssr, err := i.ssr(page)
		if err != nil {
			return err
		}

		viewData["ssr"] = ssr
	} else {
		viewData["ssr"] = nil
	}

	ts, err := i.createRootTemplate()
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html")

	err = ts.Execute(w, viewData)
	if err != nil {
		return err
	}

	return nil
}

// Location function.
func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string) {
	if i.isInertiaRequest(r) {
		w.Header().Set("X-Inertia-Location", url)
		w.WriteHeader(http.StatusConflict)
	} else if r.Method == http.MethodPost || r.Method == http.MethodPatch || r.Method == http.MethodPut {
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (i *Inertia) Back(w http.ResponseWriter, r *http.Request) {
	i.Location(w, r, r.Referer())
}

func (i *Inertia) isInertiaRequest(r *http.Request) bool {
	return r.Header.Get("X-Inertia") != ""
}

func (i *Inertia) createRootTemplate() (*template.Template, error) {
	ts := template.New(filepath.Base(i.rootTemplate)).Funcs(i.SharedFuncMap)

	if i.templateFS != nil {
		return ts.ParseFS(i.templateFS, i.rootTemplate)
	}

	return ts.ParseFiles(i.rootTemplate)
}

func (i *Inertia) ssr(page *Page) (*Ssr, error) {
	body, err := json.Marshal(page)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		strings.ReplaceAll(i.SsrURL, "/render", "")+"/render",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := i.SsrClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, ErrBadSsrStatusCode
	}

	var ssr Ssr

	err = json.NewDecoder(resp.Body).Decode(&ssr)
	if err != nil {
		return nil, err
	}

	return &ssr, nil
}

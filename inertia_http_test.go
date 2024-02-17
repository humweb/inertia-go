package inertia

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"html"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type InertiaHttpTestSuite struct {
	suite.Suite
}

//func (suite *InertiaHttpTestSuite) SetupSuite() {
//	// Setup config and ENV variables
//
//}

func (suite *InertiaHttpTestSuite) TestShare() {

	w, r := mockRequest("GET", "/users", Headers{
		"X-Inertia": "true",
	})

	i := New("", "", "")
	i.Share("title", "Page title")
	ctx := i.WithProps(r.Context(), Props{"foo": "baz", "abc": "456", "ctx": "prop"})

	err := i.Render(w, r.WithContext(ctx), "Users", Props{
		"user": map[string]interface{}{
			"name": "foo",
		},
	})

	suite.Nil(err)
	var page Page
	err = json.Unmarshal(w.Body.Bytes(), &page)
	suite.Nil(err)

	suite.Equal("Users", page.Component)

	user := page.Props["user"].(map[string]interface{})
	suite.Equal("foo", user["name"])
	suite.Equal("Page title", page.Props["title"])
	suite.Equal("baz", page.Props["foo"])
	suite.Equal("456", page.Props["abc"])
	suite.Equal("prop", page.Props["ctx"])
}

func (suite *InertiaHttpTestSuite) TestLazyProps() {

	w, r := mockRequest("GET", "/users", Headers{
		"X-Inertia":                   "true",
		"X-Inertia-Partial-Component": "Users",
		"X-Inertia-Partial-Data":      "lazy,user",
	})

	i := New("", "", "")
	i.Share("title", "Page title")
	ctx := i.WithProps(r.Context(), Props{
		"foo": "bar",
		"lazy": LazyProp(func() (any, error) {
			return "lazyprop", nil
		}),
	})

	err := i.Render(w, r.WithContext(ctx), "Users", Props{
		"user": map[string]interface{}{
			"name": "foo",
		},
	})

	suite.Nil(err)
	var page Page
	err = json.Unmarshal(w.Body.Bytes(), &page)
	suite.Nil(err)

	suite.Equal("Users", page.Component)

	user := page.Props["user"].(map[string]interface{})
	suite.Equal("foo", user["name"])
	suite.Nil(page.Props["title"])
	suite.Nil(page.Props["foo"])
	suite.Equal("lazyprop", page.Props["lazy"])
}

func (suite *InertiaHttpTestSuite) TestLazyPropsWithoutOnly() {

	w, r := mockRequest("GET", "/users", Headers{
		"X-Inertia": "true",
	})

	i := New("", "", "")
	i.Share("title", "Page title")
	ctx := i.WithProps(r.Context(), Props{
		"foo": "bar",
		"lazy": LazyProp(func() (any, error) {
			return "lazyprop", nil
		}),
	})
	ctx = i.WithProps(ctx, Props{
		"user": map[string]interface{}{
			"name": "foo",
		},
	})
	err := i.Render(w, r.WithContext(ctx), "Users", nil)

	suite.Nil(err)
	var page Page
	err = json.Unmarshal(w.Body.Bytes(), &page)
	suite.Nil(err)

	suite.Equal("Users", page.Component)

	user := page.Props["user"].(map[string]interface{})
	suite.Equal("foo", user["name"])
	suite.Equal("Page title", page.Props["title"])
	suite.Equal("bar", page.Props["foo"])
	suite.Nil(page.Props["lazy"])
}

func (suite *InertiaHttpTestSuite) TestWithProp() {

	w, r := mockRequest("GET", "/users", Headers{
		"X-Inertia": "true",
	})

	i := New("", "", "")
	ctx := i.WithProp(r.Context(), "foo", "bar")
	ctx = i.WithProp(ctx, "ctx", "baz")

	ctx = i.WithProps(ctx, Props{
		"user": map[string]interface{}{
			"name": "foo",
		},
	})
	err := i.Render(w, r.WithContext(ctx), "Users", nil)

	suite.Nil(err)
	var page Page
	err = json.Unmarshal(w.Body.Bytes(), &page)
	suite.Nil(err)

	suite.Equal("Users", page.Component)

	user := page.Props["user"].(map[string]interface{})
	suite.Equal("foo", user["name"])
	suite.Equal("baz", page.Props["ctx"])
	suite.Equal("bar", page.Props["foo"])
}

func (suite *InertiaHttpTestSuite) TestWithViewData() {

	w, r := mockRequest("GET", "/users", Headers{})

	i := New("", "./index_test.html", "")

	ctx := i.WithViewData(r.Context(), "foo", "wtf-dude")

	err := i.Render(w, r.WithContext(ctx), "Users", nil)

	suite.Nil(err)

	var page Page
	err = json.Unmarshal(w.Body.Bytes(), &page)

	suite.Contains(html.UnescapeString(w.Body.String()), "wtf-dude")

}
func (suite *InertiaHttpTestSuite) TestMiddleware() {
	i := New("", "./index_test.html", "")
	w, r := mockRequest("GET", "/users", Headers{"X-Inertia": "true"})
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	i.Middleware(testHandler).ServeHTTP(w, r)
	resp := w.Result()
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Equal("", w.Body.String())

}
func (suite *InertiaHttpTestSuite) TestMiddlewareRedirect() {
	i := New("", "./index_test.html", "2")
	w, r := mockRequest("GET", "/users", Headers{"X-Inertia": "true"})
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	i.Middleware(testHandler).ServeHTTP(w, r)
	resp := w.Result()
	suite.False(called)
	suite.Equal(http.StatusConflict, resp.StatusCode)
	suite.Equal("/users", resp.Header.Get("X-Inertia-Location"))

}
func (suite *InertiaHttpTestSuite) TestMiddlewareSkip() {
	i := New("", "./index_test.html", "")
	w, r := mockRequest("GET", "/users", Headers{})
	called := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	i.Middleware(testHandler).ServeHTTP(w, r)
	suite.True(called)
}

func TestInertiaHttpSuite(t *testing.T) {
	suite.Run(t, new(InertiaHttpTestSuite))
}

func mockRequest(method string, target string, headers Headers) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/users", nil)
	for key, val := range headers {
		r.Header.Set(key, val)
	}
	return w, r
}

type Headers = map[string]string

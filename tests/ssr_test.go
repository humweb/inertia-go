package tests

import (
	"encoding/json"
	"github.com/humweb/inertia-go"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type InertiaSsrTestSuite struct {
	suite.Suite
}

func (suite *InertiaSsrTestSuite) TestEnableSsr() {
	i := inertia.New("", "", "")
	i.EnableSsr("ssr.test")

	suite.Equal("ssr.test", i.SsrURL)
}
func (suite *InertiaSsrTestSuite) TestEnableSsrWithDefaults() {
	i := inertia.New("", "", "")
	i.EnableSsrWithDefault()
	suite.Equal("http://127.0.0.1:13714", i.SsrURL)
}
func (suite *InertiaSsrTestSuite) TestIsSsrEnabled() {
	i := inertia.New("", "", "")

	suite.False(i.IsSsrEnabled())
	i.EnableSsrWithDefault()

	suite.True(i.IsSsrEnabled())
}
func (suite *InertiaSsrTestSuite) TestDisableSsr() {
	i := inertia.New("", "", "")
	i.EnableSsrWithDefault()
	i.DisableSsr()

	suite.False(i.IsSsrEnabled())
	suite.Nil(i.SsrClient)
}

func (suite *InertiaSsrTestSuite) TestSsrRequest() {
	i := inertia.New("", "./index_test.html", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/render" {
			suite.Equal("/render", r.URL.Path)
		}
		suite.Equal("application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		ssr, _ := json.Marshal(
			inertia.Ssr{
				Head: []string{"header"},
				Body: "body text"},
		)
		w.Write(ssr)
	}))

	i.EnableSsr(server.URL)
	defer server.Close()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/user", nil)

	err := i.Render(w, r, "User", inertia.Props{
		"user": "name",
	})
	suite.Nil(err)
	suite.True(i.IsSsrEnabled())

}

func TestInertiaSsrSuite(t *testing.T) {
	suite.Run(t, new(InertiaSsrTestSuite))
}

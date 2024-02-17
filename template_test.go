package inertia

import (
	"github.com/stretchr/testify/suite"
	"html/template"
	"testing"
)

type InertiaTemplateTestSuite struct {
	suite.Suite
}

func (suite *InertiaTemplateTestSuite) TestTemplateMarshal() {
	str, err := marshal(Page{
		Component: "Users",
		Props:     Props{"username": "foobar"},
		URL:       "/users",
		Version:   "1",
	})
	suite.Nil(err)
	suite.Equal(template.JS("{\"component\":\"Users\",\"props\":{\"username\":\"foobar\"},\"url\":\"/users\",\"version\":\"1\"}"), str)
}
func (suite *InertiaTemplateTestSuite) TestTemplateMarshalErr() {
	_, err := marshal(make(chan int))
	suite.NotNil(err)
}
func (suite *InertiaTemplateTestSuite) TestTemplateRaw() {
	str, err := raw([]string{"wtf", "123"})
	suite.Nil(err)
	suite.Equal(template.HTML("wtf\n123"), str)

	str, err = raw("wtf")
	suite.Nil(err)
	suite.Equal(template.HTML("wtf"), str)
}

func (suite *InertiaTemplateTestSuite) TestTemplateRawErr() {
	_, err := raw(make(chan int))

	suite.NotNil(err)

}

func TestInertiaTemplateSuite(t *testing.T) {
	suite.Run(t, new(InertiaTemplateTestSuite))
}

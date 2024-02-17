package tests

import (
	"context"
	"errors"
	"github.com/humweb/inertia-go"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InertiaTestSuite struct {
	suite.Suite
}

func (suite *InertiaTestSuite) TestShare() {
	i := inertia.New("", "", "")
	i.Share("title", "Page title")

	title, ok := i.SharedProps["title"].(string)

	suite.True(ok)
	suite.Equal("Page title", title)
}

func (suite *InertiaTestSuite) TestShareFunc() {
	i := inertia.New("", "", "")
	i.ShareFunc("asset", func(path string) (string, error) {
		return "/" + path, nil
	})

	_, ok := i.SharedFuncMap["asset"].(func(string) (string, error))
	suite.True(ok)
	//t.Error("expected: asset func, got: empty value")
}

func (suite *InertiaTestSuite) TestWithProp() {
	ctx := context.TODO()

	i := inertia.New("", "", "")
	ctx = i.WithProp(ctx, "user", "test-user")

	contextProps, ok := ctx.Value(inertia.ContextKeyProps).(inertia.Props)
	suite.True(ok)

	user, ok := contextProps["user"].(string)
	suite.True(ok)

	suite.Equal("test-user", user)
}

func (suite *InertiaTestSuite) TestWithViewData() {
	ctx := context.TODO()

	i := inertia.New("", "", "")
	ctx = i.WithViewData(ctx, "meta", "test-meta")

	contextViewData, ok := ctx.Value(inertia.ContextKeyViewData).(inertia.Props)
	suite.True(ok)

	meta, ok := contextViewData["meta"].(string)

	ctx = i.WithViewData(ctx, "foo", "foo")

	suite.True(ok)
	suite.Equal("test-meta", meta)
}

func (suite *InertiaTestSuite) TestResolvePropsClosure() {

	val, err := inertia.ResolvePropVal(func() (any, error) {
		return "foo", nil
	})
	suite.Equal("foo", val)
	suite.Nil(err)

	val, err = inertia.ResolvePropVal(func() (any, error) {
		return nil, errors.New("nothing")
	})
	suite.Error(err)
	suite.Nil(val)

}

func (suite *InertiaTestSuite) TestResolvePropsLazy() {

	val, err := inertia.ResolvePropVal(inertia.LazyProp(func() (any, error) {
		return "foo", nil
	}))

	suite.Equal("foo", val)
	suite.Nil(err)

	val, err = inertia.ResolvePropVal(inertia.LazyProp(func() (any, error) {
		return nil, errors.New("nothing")
	}))
	suite.Error(err)
	suite.Nil(val)

}

func TestInertiaSuite(t *testing.T) {
	suite.Run(t, new(InertiaTestSuite))
}

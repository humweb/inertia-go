package inertia

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type InertiaTestSuite struct {
	suite.Suite
}

//func (suite *InertiaTestSuite) SetupSuite() {
//	// Setup config and ENV variables
//
//}

func (suite *InertiaTestSuite) TestShare() {
	i := New("", "", "")
	i.Share("title", "Inertia.js Go")

	title, ok := i.sharedProps["title"].(string)

	suite.True(ok)
	suite.Equal("Inertia.js Go", title)
}

func (suite *InertiaTestSuite) TestShareFunc() {
	i := New("", "", "")
	i.ShareFunc("asset", func(path string) (string, error) {
		return "/" + path, nil
	})

	_, ok := i.sharedFuncMap["asset"].(func(string) (string, error))
	suite.True(ok)
	//t.Error("expected: asset func, got: empty value")
}

func (suite *InertiaTestSuite) TestWithProp() {
	ctx := context.TODO()

	i := New("", "", "")
	ctx = i.WithProp(ctx, "user", "test-user")

	contextProps, ok := ctx.Value(ContextKeyProps).(Props)
	suite.True(ok)

	user, ok := contextProps["user"].(string)
	suite.True(ok)

	suite.Equal("test-user", user)
}

func (suite *InertiaTestSuite) TestWithViewData() {
	ctx := context.TODO()

	i := New("", "", "")
	ctx = i.WithViewData(ctx, "meta", "test-meta")

	contextViewData, ok := ctx.Value(ContextKeyViewData).(Props)
	suite.True(ok)

	meta, ok := contextViewData["meta"].(string)

	ctx = i.WithViewData(ctx, "foo", "foo")

	suite.True(ok)
	suite.Equal("test-meta", meta)
}

func (suite *InertiaTestSuite) TestResolvePropsClosure() {

	val, err := resolvePropVal(func() (any, error) {
		return "foo", nil
	})
	suite.Equal("foo", val)
	suite.Nil(err)

	val, err = resolvePropVal(func() (any, error) {
		return nil, errors.New("nothing")
	})
	suite.Error(err)
	suite.Nil(val)

}

func (suite *InertiaTestSuite) TestResolvePropsLazy() {

	val, err := resolvePropVal(LazyProp(func() (any, error) {
		return "foo", nil
	}))

	suite.Equal("foo", val)
	suite.Nil(err)

	val, err = resolvePropVal(LazyProp(func() (any, error) {
		return nil, errors.New("nothing")
	}))
	suite.Error(err)
	suite.Nil(val)

}

func TestInertiaSuite(t *testing.T) {
	suite.Run(t, new(InertiaTestSuite))
}

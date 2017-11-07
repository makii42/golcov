package runner

import (
	"bytes"
	"fmt"
	tt "testing"

	"github.com/golang/mock/gomock"
	osmocks "github.com/makii42/golcov/mocks/osa"
	testmocks "github.com/makii42/golcov/mocks/test"
	"github.com/makii42/golcov/test"

	"github.com/stretchr/testify/assert"
)

//
// NewTestRunner tests
//

func TesNewTestRunner(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := osmocks.NewMockOS(mc)
	test1 := testmocks.NewMockTest(mc)
	test2 := testmocks.NewMockTest(mc)
	var buf bytes.Buffer

	r, err := NewTestRunner(osa, &buf, test1, test2)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	if tr, ok := r.(*testRunner); ok {
		assert.Equal(t, osa, tr.osa)
		assert.Equal(t, &buf, tr.Out)
		assert.Equal(t, test1, tr.tests[0])
		assert.Equal(t, test2, tr.tests[1])
	} else {
		t.Fail()
	}
}

func TestCreationFailsBecauseNoGoBinary(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := osmocks.NewMockOS(mc)

	r, err := NewTestRunner(osa, nil) // reader not required, no tests...

	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, "no tests specified", err.Error())
}

//
// Run - test loop tests
//

func TestRun(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := osmocks.NewMockOS(mc)
	theTest := testmocks.NewMockTest(mc)
	outcome := testmocks.NewMockOutcome(mc)
	theTest.EXPECT().Run().Return(outcome, nil)

	tr := &testRunner{
		osa:   osa,
		tests: []test.Test{theTest},
	}

	r, err := tr.Run()

	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestRunAbortsOnHardError(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	theTest := testmocks.NewMockTest(mc)
	fakeErr := fmt.Errorf("boom")
	theTest.EXPECT().Run().Return(nil, fakeErr)
	tr := &testRunner{
		tests: []test.Test{theTest},
	}

	r, err := tr.Run()

	assert.Nil(t, r)
	assert.NotNil(t, err)
	assert.Equal(t, fakeErr, err)
}

func TestDiscoverPkgs(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	wd := "/foo/bar"
	osa := osmocks.NewMockOS(mc)
	osa.EXPECT().Getwd().Return(wd, nil)
	// todo replace with function matcher
	osa.EXPECT().Walk(wd, gomock.Any()).Return(nil)
	d := &discoverer{
		osa: osa,
	}

	pkgs, err := d.DiscoverPkgs(wd)

	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
}

type funcMatcher struct {
}

func newFuncMatcher() gomock.Matcher {
	return &funcMatcher{}
}

func (fm *funcMatcher) Matches(x interface{}) bool {
	return false
}

func (fm *funcMatcher) String() string {
	return ""
}

package runner

import (
	"bytes"
	"os"
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

const (
	goBin = "/foo/bar/go"
)

func TestCreation(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := osmocks.NewMockOS(mc)
	test1 := testmocks.NewMockTest(mc)
	test2 := testmocks.NewMockTest(mc)
	var buf bytes.Buffer

	r, err := NewTestRunner(goBin, osa, &buf, test1, test2)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	if tr, ok := r.(*testRunner); ok {
		assert.Equal(t, osa, tr.osa)
		assert.Equal(t, goBin, tr.goBinary)
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

	r, err := NewTestRunner(goBin, osa, nil) // reader not required, no tests...

	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, "no tests specified", err.Error())
}

//
// oneTest - single test execution tests
//

//
// Run - test loop tests
//

func TestRun(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin := "/bin/go"
	osa := osmocks.NewMockOS(mc)
	theTest := testmocks.NewMockTest(mc)
	outcome := testmocks.NewMockOutcome(mc)
	theTest.EXPECT().Run().Return(outcome, nil)

	tr := &testRunner{
		goBinary: goBin,
		osa:      osa,
		tests:    []test.Test{theTest},
	}

	r, err := tr.Run()

	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestDiscoverPkgs(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := osmocks.NewMockOS(mc)
	wd, err := os.Getwd()
	assert.Nil(t, err)
	tr := &testRunner{
		osa: osa,
	}
	pkgs, err := tr.DiscoverPkgs(wd)
	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
}

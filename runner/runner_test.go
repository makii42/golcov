package runner

import (
	"bytes"
	"fmt"
	tt "testing"

	"github.com/golang/mock/gomock"
	"github.com/makii42/golcov/mocks"

	"github.com/stretchr/testify/assert"
)

//
// NewTestRunner tests
//

func TestCreation(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin := "/foo/bar/go"
	osa := mocks.NewMockOS(mc)
	var buf bytes.Buffer
	osa.EXPECT().LookPath("go").Return(goBin, nil)

	r, err := NewTestRunner(osa, &buf)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	if tr, ok := r.(*testRunner); ok {
		assert.Equal(t, osa, tr.osa)
		assert.Equal(t, goBin, tr.goBinary)
		assert.Equal(t, &buf, tr.Out)
	} else {
		t.Fail()
	}
}

func TestCreationFailsBecauseNoGoBinary(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	osa := mocks.NewMockOS(mc)
	fakeErr := fmt.Errorf("boom")
	osa.EXPECT().LookPath("go").Return("", fakeErr)

	r, err := NewTestRunner(osa, nil)

	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, fakeErr, err)
}

//
// oneTest - single test execution tests
//

func TestOneTestReturnsSuccessResult(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin, tfName, pkg := "/bin/go", "/some/tmp/file.out", "./somepkg"
	testOutput := []byte("some test output")
	osa := mocks.NewMockOS(mc)
	tempFile := mocks.NewMockFile(mc)
	cmd := mocks.NewMockCommand(mc)
	osa.EXPECT().TempFile("", tempfilePrefix).Return(tempFile, nil)
	tempFile.EXPECT().Name().Return(tfName)
	osa.EXPECT().Command(goBin, "test", "-cover", "-coverprofile", tfName, "-v", pkg).Return(cmd)
	cmd.EXPECT().CombinedOutput().Return(testOutput, nil)
	tr := &testRunner{
		goBinary: goBin,
		osa:      osa,
	}
	outcome, err := tr.oneTest(pkg)

	assert.Nil(t, err)
	assert.NotNil(t, outcome.output)
	assert.Equal(t, testOutput, outcome.output)
}

func TestOneRunFailsBecauseNoTempFile(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	fakeErr := fmt.Errorf("no temp")
	osa := mocks.NewMockOS(mc)
	osa.EXPECT().TempFile("", tempfilePrefix).Return(nil, fakeErr)

	tr := &testRunner{
		osa: osa,
	}
	outcome, err := tr.oneTest("./somepkg")

	assert.Nil(t, outcome)
	assert.NotNil(t, err)
	assert.Equal(t, fakeErr, err)
}

func TestOneRunFailsBecauseCommandErrors(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin, pkg, tfName := "/bin/go", "./somepkg", "/some/tmp/file.out"
	fakeErr := fmt.Errorf("no out")
	osa := mocks.NewMockOS(mc)
	tempFile := mocks.NewMockFile(mc)
	cmd := mocks.NewMockCommand(mc)
	osa.EXPECT().TempFile("", tempfilePrefix).Return(tempFile, nil)
	tempFile.EXPECT().Name().Return(tfName)
	osa.EXPECT().Command(goBin, "test", "-cover", "-coverprofile", tfName, "-v", pkg).Return(cmd)
	cmd.EXPECT().CombinedOutput().Return(nil, fakeErr)

	tr := &testRunner{
		goBinary: goBin,
		osa:      osa,
	}
	outcome, err := tr.oneTest(pkg)

	assert.Nil(t, outcome)
	assert.NotNil(t, err)
	if testErr, ok := err.(*testError); ok {
		assert.Equal(t, fakeErr, testErr.original)
	} else {
		t.Fail()
	}
}

//
// Run - test loop tests
//

func TestRun(t *tt.T) {

}

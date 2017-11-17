package test

import (
	"fmt"
	tt "testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/makii42/golcov/mocks/osa"
	"github.com/stretchr/testify/assert"
)

func TestNewTest(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin, pkg := "/bin/foo", "./bla/"
	osa := mocks.NewMockOS(mc)

	test := NewTest(goBin, pkg, osa)

	assert.NotNil(t, test)
}

func TestRunReturnsSuccessResult(t *tt.T) {
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
	test := NewTest(goBin, pkg, osa)

	outcome, err := test.Run()

	assert.Nil(t, err)
	assert.NotNil(t, outcome.ConsoleOutput())
	assert.Equal(t, testOutput, outcome.ConsoleOutput())
}

func TestRunFailsBecauseNoTempFile(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	fakeErr := fmt.Errorf("no temp")
	osa := mocks.NewMockOS(mc)
	osa.EXPECT().TempFile("", tempfilePrefix).Return(nil, fakeErr)
	test := NewTest("/bin/go", "./whatever", osa)

	outcome, err := test.Run()

	assert.Nil(t, outcome)
	assert.NotNil(t, err)
	assert.Equal(t, fakeErr, err)
}

func TestRunFailsBecauseCommandErrors(t *tt.T) {
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
	test := NewTest(goBin, pkg, osa)

	outcome, err := test.Run()

	assert.Nil(t, outcome)
	assert.NotNil(t, err)
	if testErr, ok := err.(*TestFailure); ok {
		assert.Equal(t, fakeErr, testErr.Original)
	} else {
		t.Fail()
	}
}

package main

import (
	"bytes"
	"fmt"
	tt "testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func TestCreation(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	goBin := "/foo/bar/go"
	fos := NewMockOS(mc)
	ffs := afero.NewMemMapFs()
	var buf bytes.Buffer
	fos.EXPECT().LookPath("go").Return(goBin, nil)

	r, err := NewTestRunner(fos, ffs, &buf)

	assert.Nil(t, err)
	assert.NotNil(t, r)
	if tr, ok := r.(*testRunner); ok {
		assert.Equal(t, ffs, tr.fs)
		assert.Equal(t, fos, tr.os)
		assert.Equal(t, goBin, tr.goBinary)
		assert.Equal(t, &buf, tr.Out)
	} else {
		t.Fail()
	}
}

func TestCreationFailsBecauseNoGoBinary(t *tt.T) {
	mc := gomock.NewController(t)
	defer mc.Finish()
	fos := NewMockOS(mc)
	fakeErr := fmt.Errorf("boom")
	fos.EXPECT().LookPath("go").Return("", fakeErr)

	r, err := NewTestRunner(fos, nil, nil)

	assert.Nil(t, r)
	assert.Error(t, err)
	assert.Equal(t, fakeErr, err)
}

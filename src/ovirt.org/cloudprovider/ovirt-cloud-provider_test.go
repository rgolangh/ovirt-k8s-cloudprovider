package cloudprovider

import (
	"testing"
)
type TestReader struct {

}

func TestNewProvider(t *testing.T) {
	NewOvirtProvider(&TestReader{})
}

func (t *TestReader) Read(p []byte) (n int, err error) {
	return 0, nil
}


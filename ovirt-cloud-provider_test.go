package ovirt

import (
	"testing"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"net/http/httptest"
)

var vmsJsonData string
var expectedId = "0013e1a7-c837-4b7f-8420-6deca9486415"

type HttpHandler struct{}

func TestMain(m *testing.M) {
	if vmsJsonData == "" {
		parse, err := ioutil.ReadFile("/home/rgolan/go/src/ovirtcloudprovider/vms.json")
		if err != nil {
			panic(err)
		}
		vmsJsonData = string(parse)
	}
	m.Run()
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := string(vmsJsonData)
	io.WriteString(w, i)
}

func mockGetVms() *httptest.Server {
	server := httptest.NewServer(&HttpHandler{})
	return server
}

func TestNewProvider(t *testing.T) {
	p, err := NewOvirtProvider(ProviderConfig{})
	if err != nil || p == nil {
		t.Fatal(err)
	}
}

// TestGetInstanceId test the id returned from the api call
func TestGetInstanceId(t *testing.T) {
	// mock the api call to return a json of vms
	httpServer := mockGetVms()
	defer httpServer.Close()

	c := ProviderConfig{}
	c.Connection.Url = httpServer.URL
	provider, err := NewOvirtProvider(c)

	id, err := provider.InstanceID("ovirtNode_1")
	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Fatal(err)
	}

	if id != expectedId {
		t.Fatalf("expected id %s is no equal to %s", expectedId, id)
	}
}



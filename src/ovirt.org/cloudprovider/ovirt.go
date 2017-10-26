package cloudprovider

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"gopkg.in/gcfg.v1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/ovirt"
)

const ProviderName = "ovirt"


type OvirtNode struct {
	UUID	  string
	Name	  string
	IPAddress string
}

type ProviderConfig struct {
	Connection struct {
		Url			string	`gcfg:"url"`
		Username	string	`gcfg:"username"`
		Password 	string	`gcfg:"password"`
		Insecure	bool 	`gcfg:"insecure"`
		CAFile		string	`gcfg:"cafile"`
	}
	VmsQuery 	string "k8s"
}

type OvirtCloudProvider struct {
	VmsQuery *url.URL
}

func init() {
	cloudprovider.RegisterCloudProvider(
		ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newOvirtProvider(config)
		})
}

func NewOvirtProvider(config io.Reader) (*OvirtCloudProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("missing configuration file for ovirt cloud provider")
	}

	providerConfig := ProviderConfig{}
	providerConfig.Connection.Username = "admin@internal"

	err := gcfg.ReadInto(&providerConfig, config)
	if err != nil {
		return nil, err
	}

	vmsQuery, err := url.Parse(providerConfig.Connection.Url)
	if err != nil {
		return nil, err
	}

	vmsQuery.Path = path.Join(vmsQuery.Path, "vms")
	vmsQuery.RawQuery = url.Values{"search": {providerConfig.VmsQuery}}.Encode()

	return &OvirtCloudProvider{VmsQuery: vmsQuery}, nil

}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (*OvirtCloudProvider) Initialize(clientBuilder controller.ControllerClientBuilder) {

}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (*OvirtCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (p *OvirtCloudProvider) Instances() (cloudprovider.Instances, bool) {
	return p, true
}

func (*OvirtCloudProvider)
// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (*OvirtCloudProvider) Zones() (Zones, bool)
// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (*OvirtCloudProvider) Clusters() (Clusters, bool)
// Routes returns a routes interface along with whether the interface is supported.
func (*OvirtCloudProvider) Routes() (Routes, bool)
// ProviderName returns the cloud provider ID.
func (*OvirtCloudProvider) ProviderName() string
// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (*OvirtCloudProvider) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string)
// HasClusterID returns true if a ClusterID is required and set
func (*OvirtCloudProvider) HasClusterID() bool



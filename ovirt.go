package ovirt

import (
	"fmt"
	"io"
	"path"

	"gopkg.in/gcfg.v1"

	"net/url"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"encoding/json"
	"errors"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/ovirt"
	"k8s.io/api/core/v1"
	"net"
)

const ProviderName = "ovirt"

type OvirtNode struct {
	UUID      string
	Name      string
	IPAddress string
}

type ProviderConfig struct {
	Connection struct {
		Url      string `gcfg:"url"`
		Username string `gcfg:"username"`
		Password string `gcfg:"password"`
		Insecure bool   `gcfg:"insecure"`
		CAFile   string `gcfg:"cafile"`
	}
	Filters struct {
		VmsQuery string `gcfg:"vmsquery"`
	}
}

type Provider struct {
	VmsQuery *url.URL
}

type VM struct {
	Name 		string 		`json:"name"`
	Id   		string 		`json:"id"`
	Fqdn 		string 		`json:"fqdn"`
	Addresses []net.Addr 	'json:""'
}

type VMs struct {
	Vm []VM
}

func init() {
	cloudprovider.RegisterCloudProvider(
		ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			if config == nil {
				return nil, fmt.Errorf("missing configuration file for ovirt cloud provider")
			}
			providerConfig := ProviderConfig{}
			err := gcfg.ReadInto(&providerConfig, config)
			if err != nil {
				return nil, err
			}
			return NewOvirtProvider(providerConfig)
		})
}

func NewOvirtProvider(providerConfig ProviderConfig) (*Provider, error) {
	vmsQuery, err := url.Parse(providerConfig.Connection.Url)
	if err != nil {
		return nil, err
	}

	vmsQuery.Path = path.Join(vmsQuery.Path, "vms")
	vmsQuery.RawQuery = url.Values{"search": {providerConfig.Filters.VmsQuery}}.Encode()

	return &Provider{VmsQuery: vmsQuery}, nil

}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (*Provider) Initialize(clientBuilder controller.ControllerClientBuilder) {

}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (*Provider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (p *Provider) Instances() (cloudprovider.Instances, bool) {
	return p, true
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (*Provider) Zones() (cloudprovider.Zones, bool) {
	return nil, false

}

func (p *Provider) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
	vms, err := p.getVms()
	if err == nil {
		return nil, err
	}
	var vm *VM = &vms[string(name)]
	if vm == nil {
		return nil, fmt.Errorf(
			"VM by the name %s does not exist." +
			" The VM may have been removed, or the search query criteria needs correction",
				name)
	}
	fqdn := vm.Addresses
	if fqdn == "" {
		return nil, fmt.Errorf("Missing fqdn of instance")
	}

	return
}

func (p *Provider) InstanceID(nodeName types.NodeName) (string, error) {
	vms, err := p.getVms()
	return vms[string(nodeName)].Id, err
}

func (p *Provider) getVms() (map[string]VM, error) {
	resp, err := http.Get(p.VmsQuery.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	vms := VMs{}
	err = json.NewDecoder(resp.Body).Decode(&vms)
	if err != nil {
		return nil, err
	}
	var vmsMap = make(map[string]VM)
	for i := 0; i < len(vms.Vm); i++ {
		v := vms.Vm[i]
		vmsMap[v.Name] = v
	}

	return vmsMap, nil
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (*Provider) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (*Provider) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

// ProviderName returns the cloud provider ID.
func (*Provider) ProviderName() string {
	return ProviderName
}

// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (*Provider) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil

}

// HasClusterID returns true if a ClusterID is required and set
func (*Provider) HasClusterID() bool {
	return false
}

func (*Provider) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return errors.New("NotImplemented")
}

func (*Provider) CurrentNodeName(hostname string) (types.NodeName, error) {
	//var r types.NodeName = ""
	return types.NodeName(hostname), nil
}


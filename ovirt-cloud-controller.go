package main

import (
	"fmt"
	"os"

	"k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app/options"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/version/prometheus"        // for version metric registration
	"k8s.io/kubernetes/pkg/version/verflag"

	"github.com/spf13/pflag"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"github.com/ovirt/ovirt-k8s-cloudprovider/pkg/cloud-provider"
)

func main() {
	s := options.NewCloudControllerManagerServer()
	s.AddFlags(pflag.CommandLine)

	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	verflag.PrintAndExitIfRequested()

	_, err := cloudprovider.InitCloudProvider(cloud_provider.ProviderName, s.CloudConfigFile)

	exitOnError(err)
	if err := app.Run(s); err != nil {

	}
}
func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

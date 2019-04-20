/*
Copyright 2019 The KubeOne Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// DefaultPodSubnet defines the default subnet used by pods
	DefaultPodSubnet = "10.244.0.0/16"
	// DefaultServiceSubnet defines the default subnet used by services
	DefaultServiceSubnet = "10.96.0.0/12"
	// DefaultServiceDNS defines the default DNS domain name used by services
	DefaultServiceDNS = "cluster.local"
	// DefaultNodePortRange defines the default NodePort range
	DefaultNodePortRange = "30000-32767"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_KubeOneCluster(obj *KubeOneCluster) {
	SetDefaults_Hosts(obj)
	SetDefaults_APIEndpoints(obj)
	SetDefaults_ClusterNetwork(obj)
	SetDefaults_MachineController(obj)
}

func SetDefaults_Hosts(obj *KubeOneCluster) {
	// Set first host to be the leader
	obj.Hosts[0].IsLeader = true

	// Define a unique ID for each host
	for idx := range obj.Hosts {
		obj.Hosts[idx].ID = idx
		defaultHostConfig(&obj.Hosts[idx])
	}
}

func SetDefaults_APIEndpoints(obj *KubeOneCluster) {
	if len(obj.APIEndpoints) == 0 {
		obj.APIEndpoints = []APIEndpoint{
			{
				Host: obj.Hosts[0].PublicAddress,
			},
		}
	}
}

func SetDefaults_ClusterNetwork(obj *KubeOneCluster) {
	obj.ClusterNetwork.PodSubnet = DefaultPodSubnet
	obj.ClusterNetwork.ServiceSubnet = DefaultServiceSubnet
	obj.ClusterNetwork.ServiceDomainName = DefaultServiceDNS
	obj.ClusterNetwork.NodePortRange = DefaultNodePortRange
}

func SetDefaults_MachineController(obj *KubeOneCluster) {
	if obj.MachineController == nil {
		obj.MachineController = &MachineControllerConfig{
			Deploy: true,
		}
	}

	// If ProviderName is not None default to cloud provider and ensure user have not
	// manually provided machine-controller provider different than cloud provider.
	// If ProviderName is None, take user input or default to None.
	if obj.CloudProvider.Name != CloudProviderNameNone {
		if obj.MachineController.Provider == "" {
			obj.MachineController.Provider = obj.CloudProvider.Name
		}
	}

	// TODO(xmudrii): error
	obj.MachineController.Credentials, _ = obj.MachineController.Provider.ProviderCredentials()
}

func SetDefaults_Features(obj *KubeOneCluster) {
	if obj.Features.MetricsServer == nil {
		obj.Features.MetricsServer = &MetricsServer{
			Enable: true,
		}
	}
}

func defaultHostConfig(obj *HostConfig) {
	if len(obj.PublicAddress) == 0 && len(obj.PrivateAddress) > 0 {
		obj.PublicAddress = obj.PrivateAddress
	}
	if len(obj.PrivateAddress) == 0 && len(obj.PublicAddress) > 0 {
		obj.PrivateAddress = obj.PublicAddress
	}
	if len(obj.SSHPrivateKeyFile) == 0 && len(obj.SSHAgentSocket) == 0 {
		obj.SSHAgentSocket = "env:SSH_AUTH_SOCK"
	}
	if obj.SSHUsername == "" {
		obj.SSHUsername = "root"
	}
}
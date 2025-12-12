package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	// Cluster-scoped resource configs
	dnsCluster "github.com/Marouan-chak/provider-upjet-maas/config/cluster/dns"
	infrastructureCluster "github.com/Marouan-chak/provider-upjet-maas/config/cluster/infrastructure"
	machineCluster "github.com/Marouan-chak/provider-upjet-maas/config/cluster/machine"
	networkCluster "github.com/Marouan-chak/provider-upjet-maas/config/cluster/network"
	storageCluster "github.com/Marouan-chak/provider-upjet-maas/config/cluster/storage"

	// Namespaced resource configs
	dnsNamespaced "github.com/Marouan-chak/provider-upjet-maas/config/namespaced/dns"
	infrastructureNamespaced "github.com/Marouan-chak/provider-upjet-maas/config/namespaced/infrastructure"
	machineNamespaced "github.com/Marouan-chak/provider-upjet-maas/config/namespaced/machine"
	networkNamespaced "github.com/Marouan-chak/provider-upjet-maas/config/namespaced/network"
	storageNamespaced "github.com/Marouan-chak/provider-upjet-maas/config/namespaced/storage"
)

const (
	resourcePrefix = "maas"
	modulePath     = "github.com/Marouan-chak/provider-upjet-maas"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("maas.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		networkCluster.Configure,
		dnsCluster.Configure,
		machineCluster.Configure,
		infrastructureCluster.Configure,
		storageCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("maas.m.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		networkNamespaced.Configure,
		dnsNamespaced.Configure,
		machineNamespaced.Configure,
		infrastructureNamespaced.Configure,
		storageNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

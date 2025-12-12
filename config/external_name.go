package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// Network resources
	"maas_fabric":                   config.IdentifierFromProvider,
	"maas_vlan":                     config.IdentifierFromProvider,
	"maas_subnet":                   config.IdentifierFromProvider,
	"maas_subnet_ip_range":          config.IdentifierFromProvider,
	"maas_space":                    config.IdentifierFromProvider,
	"maas_network_interface_physical": config.IdentifierFromProvider,
	"maas_network_interface_bond":   config.IdentifierFromProvider,
	"maas_network_interface_bridge": config.IdentifierFromProvider,
	"maas_network_interface_vlan":   config.IdentifierFromProvider,
	"maas_network_interface_link":   config.IdentifierFromProvider,

	// DNS resources
	"maas_dns_domain": config.IdentifierFromProvider,
	"maas_dns_record": config.IdentifierFromProvider,

	// Machine resources
	"maas_machine":         config.IdentifierFromProvider,
	// "maas_instance": config.IdentifierFromProvider, // Disabled: complex allocate_params causes example generation issues
	"maas_vm_host":         config.IdentifierFromProvider,
	"maas_vm_host_machine": config.IdentifierFromProvider,

	// Infrastructure resources
	"maas_resource_pool": config.IdentifierFromProvider,
	"maas_tag":           config.IdentifierFromProvider,
	"maas_user":          config.IdentifierFromProvider,
	"maas_device":        config.IdentifierFromProvider,

	// Storage resources
	"maas_block_device": config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}

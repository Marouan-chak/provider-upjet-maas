package network

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ShortGroupNetwork is the short group name for network resources
const ShortGroupNetwork = "network"

// Configure configures the network group resources
func Configure(p *config.Provider) {
	// Core network resources
	p.AddResourceConfigurator("maas_fabric", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "Fabric"
	})

	p.AddResourceConfigurator("maas_vlan", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "VLAN"
		// VLAN references fabric by name
		r.References["fabric"] = config.Reference{
			TerraformName: "maas_fabric",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name", true)`,
		}
	})

	p.AddResourceConfigurator("maas_subnet", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "Subnet"
		// Subnet references VLAN - VLAN uses numeric ID, so we extract from atProvider
		r.References["vlan"] = config.Reference{
			TerraformName: "maas_vlan",
			// VLAN is identified by numeric ID in MAAS
			Extractor: `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
	})

	p.AddResourceConfigurator("maas_subnet_ip_range", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "SubnetIPRange"
		// SubnetIPRange references subnet by CIDR
		r.References["subnet"] = config.Reference{
			TerraformName: "maas_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("cidr", true)`,
		}
	})

	p.AddResourceConfigurator("maas_space", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "Space"
	})

	// Network interface resources
	// These reference machines by system_id (hostname can also work in MAAS API)
	p.AddResourceConfigurator("maas_network_interface_physical", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "InterfacePhysical"
		// Machine is identified by system_id in MAAS, which is the external-name
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
	})

	p.AddResourceConfigurator("maas_network_interface_bond", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "InterfaceBond"
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
	})

	p.AddResourceConfigurator("maas_network_interface_bridge", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "InterfaceBridge"
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
	})

	p.AddResourceConfigurator("maas_network_interface_vlan", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "InterfaceVLAN"
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
		r.References["fabric"] = config.Reference{
			TerraformName: "maas_fabric",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name", true)`,
		}
	})

	p.AddResourceConfigurator("maas_network_interface_link", func(r *config.Resource) {
		r.ShortGroup = ShortGroupNetwork
		r.Kind = "InterfaceLink"
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
		r.References["subnet"] = config.Reference{
			TerraformName: "maas_subnet",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("cidr", true)`,
		}
	})
}

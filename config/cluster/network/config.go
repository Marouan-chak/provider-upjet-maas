package network

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the network group resources
func Configure(p *ujconfig.Provider) {
	// Core network resources
	p.AddResourceConfigurator("maas_fabric", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "Fabric"
	})

	p.AddResourceConfigurator("maas_vlan", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "VLAN"
		r.References["fabric"] = ujconfig.Reference{
			TerraformName: "maas_fabric",
		}
	})

	p.AddResourceConfigurator("maas_subnet", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "Subnet"
		r.References["vlan"] = ujconfig.Reference{
			TerraformName: "maas_vlan",
		}
	})

	p.AddResourceConfigurator("maas_subnet_ip_range", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "SubnetIPRange"
		r.References["subnet"] = ujconfig.Reference{
			TerraformName: "maas_subnet",
		}
	})

	p.AddResourceConfigurator("maas_space", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "Space"
	})

	// Network interface resources
	p.AddResourceConfigurator("maas_network_interface_physical", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "InterfacePhysical"
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
	})

	p.AddResourceConfigurator("maas_network_interface_bond", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "InterfaceBond"
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
	})

	p.AddResourceConfigurator("maas_network_interface_bridge", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "InterfaceBridge"
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
	})

	p.AddResourceConfigurator("maas_network_interface_vlan", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "InterfaceVLAN"
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
		r.References["fabric"] = ujconfig.Reference{
			TerraformName: "maas_fabric",
		}
	})

	p.AddResourceConfigurator("maas_network_interface_link", func(r *ujconfig.Resource) {
		r.ShortGroup = "network"
		r.Kind = "InterfaceLink"
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
		r.References["subnet"] = ujconfig.Reference{
			TerraformName: "maas_subnet",
		}
	})
}

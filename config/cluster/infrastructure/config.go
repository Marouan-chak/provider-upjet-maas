package infrastructure

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the infrastructure group resources
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("maas_resource_pool", func(r *ujconfig.Resource) {
		r.ShortGroup = "infrastructure"
		r.Kind = "ResourcePool"
	})

	p.AddResourceConfigurator("maas_tag", func(r *ujconfig.Resource) {
		r.ShortGroup = "infrastructure"
		r.Kind = "Tag"
	})

	p.AddResourceConfigurator("maas_user", func(r *ujconfig.Resource) {
		r.ShortGroup = "infrastructure"
		r.Kind = "User"
	})

	p.AddResourceConfigurator("maas_device", func(r *ujconfig.Resource) {
		r.ShortGroup = "infrastructure"
		r.Kind = "Device"
	})
}

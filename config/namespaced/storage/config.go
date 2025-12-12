package storage

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the storage group resources
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("maas_block_device", func(r *ujconfig.Resource) {
		r.ShortGroup = "storage"
		r.Kind = "BlockDevice"
		// Block device references a machine
		r.References["machine"] = ujconfig.Reference{
			TerraformName: "maas_machine",
		}
	})
}

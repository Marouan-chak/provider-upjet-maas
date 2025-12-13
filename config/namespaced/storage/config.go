package storage

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the storage group resources
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("maas_block_device", func(r *config.Resource) {
		r.ShortGroup = "storage"
		r.Kind = "BlockDevice"
		// Block device references a machine by system_id
		r.References["machine"] = config.Reference{
			TerraformName: "maas_machine",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractResourceID()`,
		}
	})
}

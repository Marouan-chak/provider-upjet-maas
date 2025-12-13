package machine

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the machine group resources
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("maas_machine", func(r *config.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "Machine"
		// Machine can reference a resource pool by name
		// The pool field expects the pool name, not ID
		r.References["pool"] = config.Reference{
			TerraformName: "maas_resource_pool",
			// Extract the name from the referenced resource's spec
			Extractor: `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name", true)`,
		}
	})

	// NOTE: maas_instance has complex allocate_params structure that causes
	// example generation issues. Uncomment when fixed.
	// p.AddResourceConfigurator("maas_instance", func(r *config.Resource) {
	// 	r.ShortGroup = "machine"
	// 	r.Kind = "Instance"
	// })

	p.AddResourceConfigurator("maas_vm_host", func(r *config.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "VMHost"
		// VM host can reference a resource pool by name
		r.References["pool"] = config.Reference{
			TerraformName: "maas_resource_pool",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name", true)`,
		}
	})

	p.AddResourceConfigurator("maas_vm_host_machine", func(r *config.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "VMHostMachine"
		// VM host machine references a VM host by name
		r.References["vm_host"] = config.Reference{
			TerraformName: "maas_vm_host",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name", true)`,
		}
	})
}

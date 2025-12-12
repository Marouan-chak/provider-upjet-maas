package machine

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the machine group resources
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("maas_machine", func(r *ujconfig.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "Machine"
		// Machine can reference a resource pool
		r.References["pool"] = ujconfig.Reference{
			TerraformName: "maas_resource_pool",
		}
	})

	// NOTE: maas_instance has complex allocate_params structure that causes
	// example generation issues. Uncomment when fixed.
	// p.AddResourceConfigurator("maas_instance", func(r *ujconfig.Resource) {
	// 	r.ShortGroup = "machine"
	// 	r.Kind = "Instance"
	// })

	p.AddResourceConfigurator("maas_vm_host", func(r *ujconfig.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "VMHost"
		// VM host can reference a resource pool
		r.References["pool"] = ujconfig.Reference{
			TerraformName: "maas_resource_pool",
		}
	})

	p.AddResourceConfigurator("maas_vm_host_machine", func(r *ujconfig.Resource) {
		r.ShortGroup = "machine"
		r.Kind = "VMHostMachine"
		// VM host machine references a VM host
		r.References["vm_host"] = ujconfig.Reference{
			TerraformName: "maas_vm_host",
		}
	})
}

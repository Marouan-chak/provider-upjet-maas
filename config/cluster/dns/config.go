package dns

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the DNS group (domain, record)
func Configure(p *ujconfig.Provider) {
	p.AddResourceConfigurator("maas_dns_domain", func(r *ujconfig.Resource) {
		r.ShortGroup = "dns"
		r.Kind = "Domain"
	})

	p.AddResourceConfigurator("maas_dns_record", func(r *ujconfig.Resource) {
		r.ShortGroup = "dns"
		r.Kind = "Record"
		// DNS Record references a domain
		r.References["domain"] = ujconfig.Reference{
			TerraformName: "maas_dns_domain",
		}
	})
}

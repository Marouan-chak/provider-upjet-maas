// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	domain "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/dns/domain"
	record "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/dns/record"
	device "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/infrastructure/device"
	resourcepool "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/infrastructure/resourcepool"
	tag "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/infrastructure/tag"
	user "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/infrastructure/user"
	machine "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/machine/machine"
	vmhost "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/machine/vmhost"
	vmhostmachine "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/machine/vmhostmachine"
	fabric "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/fabric"
	interfacebond "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/interfacebond"
	interfacebridge "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/interfacebridge"
	interfacelink "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/interfacelink"
	interfacephysical "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/interfacephysical"
	interfacevlan "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/interfacevlan"
	space "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/space"
	subnet "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/subnet"
	subnetiprange "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/subnetiprange"
	vlan "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/network/vlan"
	providerconfig "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/providerconfig"
	blockdevice "github.com/Marouan-chak/provider-upjet-maas/internal/controller/cluster/storage/blockdevice"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domain.Setup,
		record.Setup,
		device.Setup,
		resourcepool.Setup,
		tag.Setup,
		user.Setup,
		machine.Setup,
		vmhost.Setup,
		vmhostmachine.Setup,
		fabric.Setup,
		interfacebond.Setup,
		interfacebridge.Setup,
		interfacelink.Setup,
		interfacephysical.Setup,
		interfacevlan.Setup,
		space.Setup,
		subnet.Setup,
		subnetiprange.Setup,
		vlan.Setup,
		providerconfig.Setup,
		blockdevice.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domain.SetupGated,
		record.SetupGated,
		device.SetupGated,
		resourcepool.SetupGated,
		tag.SetupGated,
		user.SetupGated,
		machine.SetupGated,
		vmhost.SetupGated,
		vmhostmachine.SetupGated,
		fabric.SetupGated,
		interfacebond.SetupGated,
		interfacebridge.SetupGated,
		interfacelink.SetupGated,
		interfacephysical.SetupGated,
		interfacevlan.SetupGated,
		space.SetupGated,
		subnet.SetupGated,
		subnetiprange.SetupGated,
		vlan.SetupGated,
		providerconfig.SetupGated,
		blockdevice.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

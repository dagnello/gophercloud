// +build acceptance networking lbaas_v2 lbaasloadbalancer

package lbaas_v2

import (
	"time"
	"testing"

	base "github.com/rackspace/gophercloud/acceptance/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const loadbalancerActiveTimeoutSeconds = 6000

func TestLoadbalancers(t *testing.T) {
	base.Setup(t)
	defer base.Teardown()

	// setup
	networkID, subnetID := SetupTopology(t)

	// create Loadbalancer
	LoadbalancerID := createLoadbalancer(t, subnetID)

	// list Loadbalancers
	listLoadbalancers(t)

	// get Loadbalancer
	getLoadbalancerWaitActive(t, LoadbalancerID)

	// update Loadbalancer
	updateLoadbalancer(t, LoadbalancerID)

	// delete Loadbalancer
	deleteLoadbalancer(t, LoadbalancerID)

	// teardown
	DeleteTopology(t, networkID)
}

func createLoadbalancer(t *testing.T, subnetID string) string {
	p, err := loadbalancers.Create(base.Client, loadbalancers.CreateOpts{
		Name:         "New_Loadbalancer",
		AdminStateUp: loadbalancers.Up,
		VipSubnetID:     subnetID,
	}).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Created Loadbalancer %s", p.ID)

	return p.ID
}

func listLoadbalancers(t *testing.T) {
	err := loadbalancers.List(base.Client, loadbalancers.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		loadbalancerList, err := loadbalancers.ExtractLoadbalancers(page)
		if err != nil {
			t.Errorf("Failed to extract Loadbalancers: %v", err)
			return false, err
		}

		for _, loadbalancer := range loadbalancerList {
			t.Logf("Listing Loadbalancer: ID [%s] Name [%s] Address [%s]",
				loadbalancer.ID, loadbalancer.Name, loadbalancer.VipAddress)
		}

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func updateLoadbalancer(t *testing.T, LoadbalancerID string) {
	_, err := loadbalancers.Update(base.Client, LoadbalancerID, loadbalancers.UpdateOpts{Name: "UpdatedName"}).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Updated Loadbalancer ID [%s]", LoadbalancerID)
}

func getLoadbalancerWaitActive(t *testing.T, LoadbalancerID string) {
	start := time.Now().Second()
	for {
		time.Sleep(1 * time.Second)

		if time.Now().Second()-start >= loadbalancerActiveTimeoutSeconds {
			t.Errorf("Loadbalancer failed to go into ACTIVE provisioning status")
			return
		}

		loadbalancer, err := loadbalancers.Get(base.Client, LoadbalancerID).Extract()
		th.AssertNoErr(t, err)
		if loadbalancer.ProvisioningStatus == "ACTIVE" {
			t.Logf("Retrieved Loadbalancer ID [%s]: OperatingStatus [%s]", loadbalancer.ID, loadbalancer.ProvisioningStatus)
		}
	}
}

func deleteLoadbalancer(t *testing.T, LoadbalancerID string) {
	res := loadbalancers.Delete(base.Client, LoadbalancerID)

	th.AssertNoErr(t, res.Err)

	t.Logf("Deleted Loadbalancer %s", LoadbalancerID)
}

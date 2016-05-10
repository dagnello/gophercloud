package listeners

import "github.com/rackspace/gophercloud"

const (
	rootPath     = "lbaas"
	resourcePath = "loadbalancers"
)

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

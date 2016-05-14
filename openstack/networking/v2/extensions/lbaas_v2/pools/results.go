package pools

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// SessionPersistence represents the session persistence feature of the load
// balancing service. It attempts to force connections or requests in the same
// session to be processed by the same member as long as it is ative. Three
// types of persistence are supported:
//
// SOURCE_IP:   With this mode, all connections originating from the same source
//              IP address, will be handled by the same member of the pool.
// HTTP_COOKIE: With this persistence mode, the load balancing function will
//              create a cookie on the first request from a client. Subsequent
//              requests containing the same cookie value will be handled by
//              the same member of the pool.
// APP_COOKIE:  With this persistence mode, the load balancing function will
//              rely on a cookie established by the backend application. All
//              requests carrying the same cookie value will be handled by the
//              same member of the pool.
type SessionPersistence struct {
	// The type of persistence mode
	Type string `mapstructure:"type" json:"type"`

	// Name of cookie if persistence mode is set appropriately
	CookieName string `mapstructure:"cookie_name" json:"cookie_name,omitempty"`
}

// Pool represents a logical set of devices, such as web servers, that you
// group together to receive and process traffic. The load balancing function
// chooses a member of the pool according to the configured load balancing
// method to handle the new requests or connections received on the VIP address.
// There is only one pool per virtual IP.
type Pool struct {
	// The load-balancer algorithm, which is round-robin, least-connections, and
	// so on. This value, which must be supported, is dependent on the provider.
	// Round-robin must be supported.
	LBMethod string `json:"lb_algorithm" mapstructure:"lb_algorithm"`

	// The protocol of the pool, which is TCP, HTTP, or HTTPS.
	Protocol string

	// Description for the pool.
	Description string

	// A list of listeners objects IDs.
	Listeners []map[string]interface{} `mapstructure:"listeners" json:"listeners"`

	// A list of member objects IDs.
	Members []map[string]interface{} `mapstructure:"members" json:"members"`

	// The IDs of associated monitors which check the health of the pool members.
	MonitorIDs []string `json:"health_monitors" mapstructure:"health_monitors"`

	// The network on which the members of the pool will be located. Only members
	// that are on this network can be added to the pool.
	SubnetID string `json:"subnet_id" mapstructure:"subnet_id"`

	// Owner of the pool. Only an administrative user can specify a tenant ID
	// other than its own.
	TenantID string `json:"tenant_id" mapstructure:"tenant_id"`

	// The administrative state of the pool, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up" mapstructure:"admin_state_up"`

	// Pool name. Does not have to be unique.
	Name string

	// List of member IDs that belong to the pool.
	// MemberIDs []string `json:"members" mapstructure:"members"`

	// The unique ID for the pool.
	ID string

	// A list of load balancer objects IDs.
	Loadbalancers []map[string]interface{} `mapstructure:"loadbalancers" json:"loadbalancers"`

	// Indicates whether connections in the same session will be processed by the
	// same pool member or not.
	Persistence SessionPersistence `mapstructure:"session_persistence" json:"session_persistence"`

	// The provider
	Provider string
}

// PoolPage is the page returned by a pager when traversing over a
// collection of pools.
type PoolPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of pools has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (p PoolPage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"pools_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// IsEmpty checks whether a PoolPage struct is empty.
func (p PoolPage) IsEmpty() (bool, error) {
	is, err := ExtractPools(p)
	spew.Dump(err)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractPools accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPools(page pagination.Page) ([]Pool, error) {
	var resp struct {
		Pools []Pool `mapstructure:"pools" json:"pools"`
	}

	err := mapstructure.Decode(page.(PoolPage).Body, &resp)

	return resp.Pools, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) Extract() (*Pool, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Pool *Pool `json:"pool"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Pool, err
}

// Member represents the application running on a backend server.
type Member struct {
	// Name of the pool.
	Name string `json:"name" mapstructure:"name"`

	// Weight of member.
	Weight int `json:"weight" mapstructure:"weight"`

	// The administrative state of the member, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up" mapstructure:"admin_state_up"`

	// Owner of the member. Only an administrative user can specify a tenant ID
	// other than its own.
	TenantID string `json:"tenant_id" mapstructure:"tenant_id"`

	// parameter value for the subnet UUID.
	SubnetID string `json:"subnet_id" mapstructure:"subnet_id"`

	// The pool to which the member belongs.
	PoolID string `json:"pool_id" mapstructure:"pool_id"`

	// The IP address of the member.
	Address string

	// The port on which the application is hosted.
	ProtocolPort int `json:"protocol_port" mapstructure:"protocol_port"`

	// The unique ID for the member.
	ID string
}

// MemberPage is the page returned by a pager when traversing over a
// collection of members in a pool.
type MemberPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of members has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (p MemberPage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"pools_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// IsEmpty checks whether a MemberPage struct is empty.
func (p MemberPage) IsEmpty() (bool, error) {
	is, err := ExtractMembers(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractMembers accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractMembers(page pagination.Page) ([]Member, error) {
	var resp struct {
		Member []Member `mapstructure:"members" json:"members"`
	}

	err := mapstructure.Decode(page.(MemberPage).Body, &resp)

	return resp.Member, err
}

type commonMemberResult struct {
	gophercloud.Result
}

// ExtractMember is a function that accepts a result and extracts a router.
func (r commonMemberResult) ExtractMember() (*Member, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Member *Member `json:"member"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Member, err
}

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) ExtractMember() (*Member, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		Member *Member `json:"member"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.Member, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// AssociateResult represents the result of an association operation.
type AssociateResult struct {
	commonResult
}

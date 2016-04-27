package common

import "strings"

// These are the various privileges a token can have.
const (
	PrivilegeRead             = 1 << iota // pretty much public data: leaderboard, scores, user profiles (without confidential stuff like email)
	PrivilegeReadConfidential             // (eventual) private messages, reports... of self
	PrivilegeWrite                        // change user information, write into confidential stuff...
	PrivilegeManageBadges                 // can change various users' badges.
	PrivilegeBetaKeys                     // can add, remove, upgrade/downgrade, make public beta keys.
	PrivilegeManageSettings               // maintainance, set registrations, global alerts, bancho settings
	PrivilegeViewUserAdvanced             // can see user email, and perhaps warnings in the future, basically.
	PrivilegeManageUser                   // can change user email, allowed status, userpage, rank, username...
	PrivilegeManageRoles                  // translates as admin, as they can basically assign roles to anyone, even themselves
	PrivilegeManageAPIKeys                // admin permission to manage user permission, not only self permissions. Only ever do this if you completely trust the application, because this essentially means to put the entire ripple database in the hands of a (potentially evil?) application.
	PrivilegeBlog                         // can do pretty much anything to the blog, and the documentation.
	PrivilegeAPIMeta                      // can do /meta API calls. basically means they can restart the API server.
	PrivilegeBeatmap                      // rank/unrank beatmaps. also BAT when implemented
)

// Privileges is a bitwise enum of the privileges of an user's API key.
type Privileges uint64

// HasPrivilegeRead returns whether the Read privilege is included in the privileges.
func (p Privileges) HasPrivilegeRead() bool {
	return p&PrivilegeRead != 0
}

// HasPrivilegeReadConfidential returns whether the ReadConfidential privilege is included in the privileges.
func (p Privileges) HasPrivilegeReadConfidential() bool {
	return p&PrivilegeReadConfidential != 0
}

// HasPrivilegeWrite returns whether the Write privilege is included in the privileges.
func (p Privileges) HasPrivilegeWrite() bool {
	return p&PrivilegeWrite != 0
}

// HasPrivilegeManageBadges returns whether the ManageBadges privilege is included in the privileges.
func (p Privileges) HasPrivilegeManageBadges() bool {
	return p&PrivilegeManageBadges != 0
}

// HasPrivilegeBetaKeys returns whether the BetaKeys privilege is included in the privileges.
func (p Privileges) HasPrivilegeBetaKeys() bool {
	return p&PrivilegeBetaKeys != 0
}

// HasPrivilegeManageSettings returns whether the ManageSettings privilege is included in the privileges.
func (p Privileges) HasPrivilegeManageSettings() bool {
	return p&PrivilegeManageSettings != 0
}

// HasPrivilegeViewUserAdvanced returns whether the ViewUserAdvanced privilege is included in the privileges.
func (p Privileges) HasPrivilegeViewUserAdvanced() bool {
	return p&PrivilegeViewUserAdvanced != 0
}

// HasPrivilegeManageUser returns whether the ManageUser privilege is included in the privileges.
func (p Privileges) HasPrivilegeManageUser() bool {
	return p&PrivilegeManageUser != 0
}

// HasPrivilegeManageRoles returns whether the ManageRoles privilege is included in the privileges.
func (p Privileges) HasPrivilegeManageRoles() bool {
	return p&PrivilegeManageRoles != 0
}

// HasPrivilegeManageAPIKeys returns whether the ManageAPIKeys privilege is included in the privileges.
func (p Privileges) HasPrivilegeManageAPIKeys() bool {
	return p&PrivilegeManageAPIKeys != 0
}

// HasPrivilegeBlog returns whether the Blog privilege is included in the privileges.
func (p Privileges) HasPrivilegeBlog() bool {
	return p&PrivilegeBlog != 0
}

// HasPrivilegeAPIMeta returns whether the APIMeta privilege is included in the privileges.
func (p Privileges) HasPrivilegeAPIMeta() bool {
	return p&PrivilegeAPIMeta != 0
}

// HasPrivilegeBeatmap returns whether the Beatmap privilege is included in the privileges.
func (p Privileges) HasPrivilegeBeatmap() bool {
	return p&PrivilegeBeatmap != 0
}

var privilegeString = [...]string{
	"Read",
	"ReadConfidential",
	"Write",
	"ManageBadges",
	"BetaKeys",
	"ManageSettings",
	"ViewUserAdvanced",
	"ManageUser",
	"ManageRoles",
	"ManageAPIKeys",
	"Blog",
	"APIMeta",
	"Beatmap",
}

func (p Privileges) String() string {
	var pvs []string
	for i, v := range privilegeString {
		if int(p)&(1<<uint(i)) != 0 {
			pvs = append(pvs, v)
		}
	}
	return strings.Join(pvs, ", ")
}

var privilegeMustBe = [...]int{
	1,
	1,
	1,
	3,
	3,
	4,
	4,
	4,
	4,
	4,
	3,
	4,
	4,
}

// CanOnly removes any privilege that the user has requested to have, but cannot have due to their rank.
func (p Privileges) CanOnly(rank int) Privileges {
	newPrivilege := 0
	for i, v := range privilegeMustBe {
		wants := p&1 == 1
		can := rank >= v
		if wants && can {
			newPrivilege |= 1 << uint(i)
		}
		p >>= 1
	}
	return Privileges(newPrivilege)
}

package v1

import (
	"github.com/osuripple/api/common"
)

type privilegesData struct {
	Read             bool `json:"read"`
	ReadConfidential bool `json:"read_confidential"`
	Write            bool `json:"write"`
	ManageBadges     bool `json:"manage_badges"`
	BetaKeys         bool `json:"beta_keys"`
	ManageSettings   bool `json:"manage_settings"`
	ViewUserAdvanced bool `json:"view_user_advanced"`
	ManageUser       bool `json:"manage_user"`
	ManageRoles      bool `json:"manage_roles"`
	ManageAPIKeys    bool `json:"manage_api_keys"`
	Blog             bool `json:"blog"`
	APIMeta          bool `json:"api_meta"`
}

// PrivilegesGET returns an explaination for the privileges, telling the client what they can do with this token.
func PrivilegesGET(md common.MethodData) (r common.Response) {
	// This code sucks.
	r.Code = 200
	r.Data = privilegesData{
		Read:             md.User.Privileges.HasPrivilegeRead(),
		ReadConfidential: md.User.Privileges.HasPrivilegeReadConfidential(),
		Write:            md.User.Privileges.HasPrivilegeWrite(),
		ManageBadges:     md.User.Privileges.HasPrivilegeManageBadges(),
		BetaKeys:         md.User.Privileges.HasPrivilegeBetaKeys(),
		ManageSettings:   md.User.Privileges.HasPrivilegeManageSettings(),
		ViewUserAdvanced: md.User.Privileges.HasPrivilegeViewUserAdvanced(),
		ManageUser:       md.User.Privileges.HasPrivilegeManageUser(),
		ManageRoles:      md.User.Privileges.HasPrivilegeManageRoles(),
		ManageAPIKeys:    md.User.Privileges.HasPrivilegeManageAPIKeys(),
		Blog:             md.User.Privileges.HasPrivilegeBlog(),
		APIMeta:          md.User.Privileges.HasPrivilegeAPIMeta(),
	}
	return
}

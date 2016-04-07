package v1

import (
	"github.com/osuripple/api/common"
)

type privilegesData struct {
	PrivilegeRead             bool `json:"read"`
	PrivilegeReadConfidential bool `json:"read_confidential"`
	PrivilegeWrite            bool `json:"write"`
	PrivilegeManageBadges     bool `json:"manage_badges"`
	PrivilegeBetaKeys         bool `json:"beta_keys"`
	PrivilegeManageSettings   bool `json:"manage_settings"`
	PrivilegeViewUserAdvanced bool `json:"view_user_advanced"`
	PrivilegeManageUser       bool `json:"manage_user"`
	PrivilegeManageRoles      bool `json:"manage_roles"`
	PrivilegeManageAPIKeys    bool `json:"manage_api_keys"`
	PrivilegeBlog             bool `json:"blog"`
	PrivilegeAPIMeta          bool `json:"api_meta"`
}

// PrivilegesGET returns an explaination for the privileges, telling the client what they can do with this token.
func PrivilegesGET(md common.MethodData) (r common.Response) {
	// This code sucks.
	r.Code = 200
	r.Data = privilegesData{
		PrivilegeRead:             md.User.Privileges.HasPrivilegeRead(),
		PrivilegeReadConfidential: md.User.Privileges.HasPrivilegeReadConfidential(),
		PrivilegeWrite:            md.User.Privileges.HasPrivilegeWrite(),
		PrivilegeManageBadges:     md.User.Privileges.HasPrivilegeManageBadges(),
		PrivilegeBetaKeys:         md.User.Privileges.HasPrivilegeBetaKeys(),
		PrivilegeManageSettings:   md.User.Privileges.HasPrivilegeManageSettings(),
		PrivilegeViewUserAdvanced: md.User.Privileges.HasPrivilegeViewUserAdvanced(),
		PrivilegeManageUser:       md.User.Privileges.HasPrivilegeManageUser(),
		PrivilegeManageRoles:      md.User.Privileges.HasPrivilegeManageRoles(),
		PrivilegeManageAPIKeys:    md.User.Privileges.HasPrivilegeManageAPIKeys(),
		PrivilegeBlog:             md.User.Privileges.HasPrivilegeBlog(),
		PrivilegeAPIMeta:          md.User.Privileges.HasPrivilegeAPIMeta(),
	}
	return
}

package v1

import (
	"git.zxq.co/ripple/rippleapi/common"
)

type privilegesData struct {
	common.ResponseBase
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
	Beatmap          bool `json:"beatmap"`
}

// PrivilegesGET returns an explaination for the privileges, telling the client what they can do with this token.
func PrivilegesGET(md common.MethodData) common.CodeMessager {
	r := privilegesData{}
	r.Code = 200
	// This code sucks.
	r.Read = md.User.Privileges.HasPrivilegeRead()
	r.ReadConfidential = md.User.Privileges.HasPrivilegeReadConfidential()
	r.Write = md.User.Privileges.HasPrivilegeWrite()
	r.ManageBadges = md.User.Privileges.HasPrivilegeManageBadges()
	r.BetaKeys = md.User.Privileges.HasPrivilegeBetaKeys()
	r.ManageSettings = md.User.Privileges.HasPrivilegeManageSettings()
	r.ViewUserAdvanced = md.User.Privileges.HasPrivilegeViewUserAdvanced()
	r.ManageUser = md.User.Privileges.HasPrivilegeManageUser()
	r.ManageRoles = md.User.Privileges.HasPrivilegeManageRoles()
	r.ManageAPIKeys = md.User.Privileges.HasPrivilegeManageAPIKeys()
	r.Blog = md.User.Privileges.HasPrivilegeBlog()
	r.APIMeta = md.User.Privileges.HasPrivilegeAPIMeta()
	r.Beatmap = md.User.Privileges.HasPrivilegeBeatmap()
	return r
}

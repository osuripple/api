package v1

import (
	"strings"

	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"zxq.co/ripple/playstyle"
	"zxq.co/ripple/rippleapi/common"
	semanticiconsugc "zxq.co/ripple/semantic-icons-ugc"
)

type donorInfoResponse struct {
	common.ResponseBase
	HasDonor   bool                 `json:"has_donor"`
	Expiration common.UnixTimestamp `json:"expiration"`
}

// UsersSelfDonorInfoGET returns information about the users' donor status
func UsersSelfDonorInfoGET(md common.MethodData) common.CodeMessager {
	var r donorInfoResponse
	var privileges uint64
	err := md.DB.QueryRow("SELECT privileges, donor_expire FROM users WHERE id = ?", md.ID()).
		Scan(&privileges, &r.Expiration)
	if err != nil {
		md.Err(err)
		return Err500
	}
	r.HasDonor = common.UserPrivileges(privileges)&common.UserPrivilegeDonor > 0
	r.Code = 200
	return r
}

type favouriteModeResponse struct {
	common.ResponseBase
	FavouriteMode  int `json:"favourite_mode"`
	FavouriteRelax int `json:"favourite_relax"`
}

// UsersSelfFavouriteModeGET gets the current user's favourite mode
func UsersSelfFavouriteModeGET(md common.MethodData) common.CodeMessager {
	var f favouriteModeResponse
	f.Code = 200
	if md.ID() == 0 {
		return f
	}
	err := md.DB.QueryRow("SELECT users_stats.favourite_mode, users_stats.favourite_relax FROM users_stats WHERE id = ?", md.ID()).
		Scan(&f.FavouriteMode, &f.FavouriteRelax)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return f
}

type userSettingsData struct {
	UsernameAKA        *string `json:"username_aka"`
	FavouriteMode      *int    `json:"favourite_mode"`
	FavouriteRelaxMode *int    `json:"favourite_relax"`
	CustomBadge        struct {
		singleBadge
		Show *bool `json:"show"`
	} `json:"custom_badge"`
	PlayStyle *int `json:"play_style"`
}

// UsersSelfSettingsPOST allows to modify information about the current user.
func UsersSelfSettingsPOST(md common.MethodData) common.CodeMessager {
	var d userSettingsData
	md.Unmarshal(&d)

	// input sanitisation
	*d.UsernameAKA = common.SanitiseString(*d.UsernameAKA)
	if md.User.UserPrivileges&common.UserPrivilegeDonor > 0 {
		d.CustomBadge.Name = common.SanitiseString(d.CustomBadge.Name)
		emoji, err := emoji.LookupEmoji(d.CustomBadge.Icon)
		if err != nil {
			return common.SimpleResponse(400, "Invalid emoji")
		}
		d.CustomBadge.Icon = emoji.Value
	} else {
		d.CustomBadge.singleBadge = singleBadge{}
		d.CustomBadge.Show = nil
	}
	d.FavouriteMode = intPtrIn(0, d.FavouriteMode, 3)
	if d.FavouriteMode != nil && *d.FavouriteMode == 3 {
		v := 0
		d.FavouriteRelaxMode = &v
	}
	d.FavouriteRelaxMode = intPtrIn(0, d.FavouriteRelaxMode, 1)

	q := new(common.UpdateQuery).
		Add("s.username_aka", d.UsernameAKA).
		Add("s.favourite_mode", d.FavouriteMode).
		Add("s.custom_badge_name", d.CustomBadge.Name).
		Add("s.custom_badge_icon", d.CustomBadge.Icon).
		Add("s.show_custom_badge", d.CustomBadge.Show).
		Add("s.play_style", *d.PlayStyle&^(playstyle.Spoon|playstyle.LeapMotion|playstyle.OculusRift|playstyle.Dick|playstyle.Eggplant)).
		Add("s.favourite_relax", d.FavouriteRelaxMode)
	_, err := md.DB.Exec("UPDATE users u, users_stats s SET "+q.Fields()+" WHERE s.id = u.id AND u.id = ?", append(q.Parameters, md.ID())...)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return UsersSelfSettingsGET(md)
}

func sanitiseIconName(s string) string {
	classes := strings.Split(s, " ")
	n := make([]string, 0, len(classes))
	for _, c := range classes {
		if !in(c, n) && in(c, semanticiconsugc.SaneIcons) {
			n = append(n, c)
		}
	}
	return strings.Join(n, " ")
}

func in(a string, b []string) bool {
	for _, x := range b {
		if x == a {
			return true
		}
	}
	return false
}

type userSettingsResponse struct {
	common.ResponseBase
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Flags    uint   `json:"flags"`
	userSettingsData
}

// UsersSelfSettingsGET allows to get "sensitive" information about the current user.
func UsersSelfSettingsGET(md common.MethodData) common.CodeMessager {
	var r userSettingsResponse
	var ccb bool
	r.Code = 200
	err := md.DB.QueryRow(`
SELECT
	u.id, u.username,
	u.email, s.username_aka, s.favourite_mode,
	s.favourite_relax,
	s.show_custom_badge, s.custom_badge_icon,
	s.custom_badge_name, s.can_custom_badge,
	s.play_style, u.flags
FROM users u
LEFT JOIN users_stats s ON u.id = s.id
WHERE u.id = ?`, md.ID()).Scan(
		&r.ID, &r.Username,
		&r.Email, &r.UsernameAKA, &r.FavouriteMode,
		&r.FavouriteRelaxMode,
		&r.CustomBadge.Show, &r.CustomBadge.Icon,
		&r.CustomBadge.Name, &ccb,
		&r.PlayStyle, &r.Flags,
	)
	if err != nil {
		md.Err(err)
		return Err500
	}
	if !ccb {
		r.CustomBadge = struct {
			singleBadge
			Show *bool `json:"show"`
		}{}
	}
	return r
}

const (
	overwritePP    = 0
	overwriteScore = 1

	modeClassic = 0
	modeRelax   = 1

	autoLastOff          = 0
	autoLastMessage      = 1
	autoLastNotification = 2

	displayScore = 0
	displayPP    = 1
)

var validOverwrites = []int{overwritePP, overwriteScore}
var validModes = []int{modeClassic, modeRelax}
var validAutoLast = []int{autoLastOff, autoLastMessage, autoLastNotification}
var validDisplayModes = []int{displayScore, displayPP}

type scoreOverwrite struct {
	Std   *int `json:"std"`
	Taiko *int `json:"taiko"`
	Ctb   *int `json:"ctb"`
	Mania *int `json:"mania"`
}

type scoreboardDisplay struct {
	Classic *int `json:"classic"`
	Relax   *int `json:"relax"`
}

type autoLast struct {
	Classic *int `json:"classic"`
	Relax   *int `json:"relax"`
}

type scoreboard struct {
	Mode    *int               `json:"mode"`
	Display *scoreboardDisplay `json:"display"`
}

type scoreboardData struct {
	common.ResponseBase
	Scoreboard *scoreboard     `json:"scoreboard"`
	Overwrite  *scoreOverwrite `json:"overwrite"`
	AutoLast   *autoLast       `json:"auto_last"`
}

type scoreboardResponse struct {
	common.ResponseBase
	scoreboardData
}

// UserSelfScoreboardGET returns the  in-game scoreboard perferences
// for the current uesr
func UserSelfScoreboardGET(md common.MethodData) common.CodeMessager {
	var r scoreboardResponse
	r.Code = 200
	r.Scoreboard = &scoreboard{Display: &scoreboardDisplay{}}
	r.AutoLast = &autoLast{}
	r.Overwrite = &scoreOverwrite{}
	err := md.DB.QueryRow(`
SELECT
	u.is_relax,
	score_overwrite_std, score_overwrite_taiko,
	score_overwrite_ctb, score_overwrite_mania,
	scoreboard_display_classic, scoreboard_display_relax,
	auto_last_classic, auto_last_relax
FROM users_preferences JOIN users AS u USING(id)
WHERE id = ? LIMIT 1`, md.ID()).Scan(
		&r.Scoreboard.Mode,
		&r.Overwrite.Std, &r.Overwrite.Taiko,
		&r.Overwrite.Ctb, &r.Overwrite.Mania,
		&r.Scoreboard.Display.Classic, &r.Scoreboard.Display.Relax,
		&r.AutoLast.Classic, &r.AutoLast.Relax,
	)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return r
}

// UserSelfScoreboardPOST allows users to change their in-game
// scoreboard preferences
func UserSelfScoreboardPOST(md common.MethodData) common.CodeMessager {
	var d scoreboardData
	err := md.Unmarshal(&d)
	if err != nil {
		return ErrBadJSON
	}
	q := new(common.UpdateQuery)

	type scoreboardDataField struct {
		value         *int
		allowedValues []int
		column        string
	}
	for _, field := range []scoreboardDataField{
		{d.Scoreboard.Mode, []int{modeClassic, modeRelax}, "u.is_relax"},
		{d.Scoreboard.Display.Classic, validDisplayModes, "scoreboard_display_classic"},
		{d.Scoreboard.Display.Relax, validDisplayModes, "scoreboard_display_relax"},
		{d.Overwrite.Std, validOverwrites, "score_overwrite_std"},
		{d.Overwrite.Taiko, validOverwrites, "score_overwrite_taiko"},
		{d.Overwrite.Ctb, validOverwrites, "score_overwrite_ctb"},
		{d.Overwrite.Mania, validOverwrites, "score_overwrite_mania"},
		{d.AutoLast.Classic, validAutoLast, "auto_last_classic"},
		{d.AutoLast.Relax, validAutoLast, "auto_last_relax"},
	} {
		if field.value == nil {
			continue
		}
		if !contains(*field.value, field.allowedValues) {
			return ErrBadField(field.column)
		}
		q.Add(field.column, field.value)
	}
	_, err = md.DB.Exec(
		`UPDATE users_preferences, users AS u
		SET `+q.Fields()+`
		WHERE users_preferences.id = ? AND u.id = ?
		LIMIT 1`,
		append(q.Parameters, md.ID(), md.ID())...,
	)
	if err != nil {
		md.Err(err)
		return Err500
	}
	return UserSelfScoreboardGET(md)
}

func contains(needle int, haystack []int) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func intPtrIn(x int, y *int, z int) *int {
	if y == nil {
		return nil
	}
	if *y > z {
		return nil
	}
	if *y < x {
		return nil
	}
	return y
}

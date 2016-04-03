package common

// Token Is an API token.
type Token struct {
	Value      string
	UserID     int
	Privileges Privileges
}

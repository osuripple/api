package common

// Token Is an API token.
type Token struct {
	ID         int
	Value      string
	UserID     int
	Privileges Privileges
}

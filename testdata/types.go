package testdata

// User is a data type for user
type User struct {
	// an exported field
	Exported string
	// this field should not be visible
	unexported string
	// name for user. required
	Named            string `json:"name"`
	NamedNotRequired string `json:"nickname,omitempty"`
	NotRequired      string `json:",omitempty"`
	Hidden           string `json:"-"`
	Status           UserStatus
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

func Foo() {
}

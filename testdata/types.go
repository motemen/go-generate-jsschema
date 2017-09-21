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
	Tags             []string
	Items            []*Item
	*Embedded
}

type Embedded struct {
	ID uint64
}

type UserMap map[string]User

type UserStatus string

type Item struct {
	Name  string
	Count int
}

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

func Foo() {
}

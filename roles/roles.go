package roles

// Role is the role type of the user
type Role string

// CanCreateHole checks if a
func (r Role) CanCreateHole() bool {
	return r == "admin"
}

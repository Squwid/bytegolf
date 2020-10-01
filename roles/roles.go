package roles

// Role is the role type of the user
type Role string

// CanCreateHole checks if a
func (r Role) CanCreateHole() bool {
	return r.isAdmin()
}

func (r Role) CanListAdminHoles() bool {
	return r.isAdmin()
}

func (r Role) isAdmin() bool {
	return r == "admin"
}

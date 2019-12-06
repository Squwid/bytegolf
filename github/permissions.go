package github

// IsGamemaster checks to see if a user is a game master which means they can create and manage
// questions but nothing else
func (u *User) IsGamemaster() bool {
	// this is a temp id and will be changed before commiting so it wont work
	return u.BGID == "9abf282b-9d05-4b11-b4d4-3f8a8af9ea2f"
}

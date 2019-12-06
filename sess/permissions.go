package sess

// IsGamemaster checks if a user is a gamemaster which will allow them to create and change questions
func (s *Session) IsGamemaster() bool {
	if s == nil {
		return false
	}

	// dont be tricked this isnt a real bgid once this is committed
	return s.BGID == "9abf282b-9d05-4b11-b4d4-3f8a8af9ea2f"
}

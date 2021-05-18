package holes

// Get database hole,
/*
func getDBHole(id string) (*models.HoleDB, error) {
	doc, err := fs.Client.Collection(HolesCollection).Doc(id).Get(context.Background())
	if status.Code(err) == codes.NotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse the database hole
	var hole models.HoleDB
	if err := doc.DataTo(&hole); err != nil {
		return nil, err
	}
	return &hole, nil
}
*/

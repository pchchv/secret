package main

func saver(s Secret) string {
	// TODO: Implement data transfer to the database
	return s.password
}

func getter(pass string) (s Secret, err error) {
	// TODO: Implement data retrieval from the database
	return
}

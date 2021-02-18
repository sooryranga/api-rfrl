package user

// createUser use case to create a new user
func (h *Handler) createUser(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) (*User, error) {
	user := newUser(
		firstName,
		lastName,
		about,
		email,
		photo,
	)
	return CreateUser(h.db, user)
}

// updateUser use case to update a new user
func (h *Handler) updateUser(
	id int,
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) (*User, error) {
	user := newUser(
		firstName,
		lastName,
		about,
		email,
		photo,
	)

	return UpdateUser(h.db, id, user)
}

func (h *Handler) getUser(id string) (*User, error) {
	return GetUserFromID(h.db, id)
}

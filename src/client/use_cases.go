package client

// createClient use case to create a new client
func (h *Handler) createClient(
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) (*Client, error) {
	client := NewClient(
		firstName,
		lastName,
		about,
		email,
		photo,
	)
	return CreateClient(h.db, client)
}

// updateClient use case to update a new client
func (h *Handler) updateClient(
	id string,
	firstName string,
	lastName string,
	about string,
	email string,
	photo string,
) (*Client, error) {
	client := NewClient(
		firstName,
		lastName,
		about,
		email,
		photo,
	)

	return UpdateClient(h.db, id, client)
}

func (h *Handler) getClient(id string) (*Client, error) {
	return GetClientFromID(h.db, id)
}

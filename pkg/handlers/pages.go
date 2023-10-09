package handlers

type indexPage struct {
	Metadata   metadata
	IsLoggedIn bool
	User       user
	Rooms      []room
}

type metadata struct {
	Title string
}

type user struct {
	Username string
}

type room struct {
	ID   string
	Name string
}

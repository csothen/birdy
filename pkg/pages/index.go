package pages

type IndexPage struct {
	Head  HeadData
	User  *UserData
	Rooms []RoomData
}

type ErrorPage struct {
	Code  uint
	Error error
}

type HeadData struct {
	Title string
}

type UserData struct {
	Username string
}

type RoomData struct {
	ID   string
	Name string
}

func (i *IndexPage) Template() string {
	return "base"
}

func (e *ErrorPage) Template() string {
	return "error"
}

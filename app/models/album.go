package models

//type omit *struct {}

type Album struct {
	Name string `json:"name"`
  	Username string `json:"username"`
}

type PublicAlbum struct {
  *Album
}

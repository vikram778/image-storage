package album

type PostAlbumResponse struct {
	AlbumTittle string `json:"album_tittle"`
	Message     string `json:"message"`
}

type DeleteAlbumResponse struct {
	AlbumTittle string `json:"album_tittle"`
	Message     string `json:"message"`
}

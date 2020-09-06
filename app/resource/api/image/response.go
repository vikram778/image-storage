package image

type PostImageResponse struct {
	AlbumTittle string `json:"album_tittle"`
	ImageName   string `json:"image_name"`
	ImageID     int64  `json:"image_id"`
	ImageCount  int64  `json:"image_count"`
}

type DeleteResponse struct {
	AlbumTittle string `json:"album_tittle"`
	ImagePath   string `json:"image_path"`
	Message     string `json:"message"`
}

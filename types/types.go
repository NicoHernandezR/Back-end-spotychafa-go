package types

type Mp3Store interface {
	GetMp3ByID(id int) (*Mp3, error)
	InsertMp3(Mp3) error
}

type Mp3 struct {
	ID      interface{} `json:"id" validate:"required"`
	Title   string      `json:"title" validate:"required"`
	Artist  string      `json:"artist" validate:"required"`
	Mp3File string      `json:"mp3File" validate:"required"`
}

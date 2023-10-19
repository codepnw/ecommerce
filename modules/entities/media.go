package entities

type Image struct {
	Id       string `db:"id" json:"id"`
	FileName string `db:"file_name" json:"file_name"`
	Url      string `db:"url" json:"url"`
}

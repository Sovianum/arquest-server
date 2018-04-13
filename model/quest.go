package model

type Quest struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Rating      float32 `json:"rating"`
	DataPath    string  `json:"data_path,omitempty"`
}

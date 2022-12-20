package dto

type CreateResp struct {
	Creator string `json:"creator"`
	Date    string `json:"date"`
}
type EventResp struct {
	EventType string `json:"eventType"`
	FileName  string `json:"fileName"`
	CoverImg  string `json:"coverImg"`
	FileHash  string `json:"fileHash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Price     string `json:"price"`
	State     string `json:"state"`
	Date      string `json:"date"`
}

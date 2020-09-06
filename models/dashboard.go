package models

// Dashboard struct
type Dashboard struct {
	ID       int64 `json:"id"`
	Request  int64 `json:"request"`
	Staging  int64 `json:"staging"`
	Queue    int64 `json:"queue"`
	Complete int64 `json:"complete"`
	Failed   int64 `json:"failed"`
}

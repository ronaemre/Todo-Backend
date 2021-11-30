package models

type Todo struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Details string `json:"details"`
}

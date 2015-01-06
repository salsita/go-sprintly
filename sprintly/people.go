package sprintly

type User struct {
	Id        int    `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Admin     bool   `json:"admin,omitempty"`
}

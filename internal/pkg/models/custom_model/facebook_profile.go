package custommodel

type FacebookUserProfile struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`
	ID         string `json:"id"`
}

// Struct for the JSON structure
type Message struct {
	CreatedTime string `json:"created_time"`
	From        struct {
		Email string `json:"email"`
		ID    string `json:"id"`
		Name  string `json:"name"`
	} `json:"from"`
	ID      string `json:"id"`
	Message string `json:"message"`
	To      struct {
		Data []struct {
			Email string `json:"email"`
			ID    string `json:"id"`
			Name  string `json:"name"`
		} `json:"data"`
	} `json:"to"`
}

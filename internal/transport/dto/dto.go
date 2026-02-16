package dto

type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

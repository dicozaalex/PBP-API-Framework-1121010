package controllers

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType int    `json:"UserType"`
}

type UserResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    User   `json:"data"`
}
type UsersResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    []User `json:"data"`
}
type ErrorResponse struct {
	Message string `json:"error"`
	Status  int    `json:"status"`
}

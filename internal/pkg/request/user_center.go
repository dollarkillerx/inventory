package request

type UserLogin struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

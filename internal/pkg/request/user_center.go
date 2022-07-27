package request

type UserLogin struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserInfo struct {
	Account string `json:"account" binding:"required"`
}

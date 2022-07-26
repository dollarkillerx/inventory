package response

type UserLogin struct {
	JWT        string `json:"jwt"`
	Storehouse string `json:"storehouse"`
}

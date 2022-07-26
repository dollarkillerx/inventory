package request

type AuthJWT struct {
	Storehouse string `json:"storehouse"` // 倉庫
	Account    string `json:"account"`
}

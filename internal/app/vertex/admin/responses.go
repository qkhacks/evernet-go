package admin

type TokenResponse struct {
	Token string `json:"token"`
}

type AdditionResponse struct {
	Password string `json:"password"`
	Admin    *Admin `json:"admin"`
}

type PasswordResponse struct {
	Password string `json:"password"`
}

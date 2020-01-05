package models

//RefreshTokenStore exported
type RefreshTokenStore interface {
	RefreshTokenVerify(string, string, int) (*RefreshToken, error)
	RefreshTokenAdd(string, string, int) (int, error)
	RefreshTokenUpdateByID(int, string, string) (int, error)
	RefreshTokenDeleteByID(int) error
}

//RefreshToken exported
type RefreshToken struct {
	ID         int    `json:"id:"`
	Identifier string `json:"identifier,omitempty"`
	Token      string `json:"refreshToken"`
	User       User   `json:"user"`
	Created    string `json:"created"`
	Updated    string `json:"updated"`
}

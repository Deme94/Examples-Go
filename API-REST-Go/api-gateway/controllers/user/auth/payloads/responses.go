package payloads

import "time"

type GetAllResponse struct {
	Users []*GetResponse `json:"users"`
}
type GetResponse struct {
	ID                 int       `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Nick               string    `json:"nick"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address"`
	LastPasswordChange time.Time `json:"last_password_change"`
	VerifiedMail       bool      `json:"verified_mail"`
	VerifiedPhone      bool      `json:"verified_phone"`
}

type LoginResponse struct {
	ID        int        `json:"user_id,omitempty"`
	Token     string     `json:"token,omitempty"`
	BanExpire *time.Time `json:"ban_expire,omitempty"`
}

type OkResponse struct {
	OK bool `json:"ok"`
}

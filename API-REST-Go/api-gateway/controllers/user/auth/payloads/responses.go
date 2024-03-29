package payloads

import "time"

type GetResponse struct {
	CreatedAt          *time.Time `json:"created_at"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	Nick               string     `json:"nick"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	Phone              string     `json:"phone"`
	Address            string     `json:"address"`
	LastPasswordChange *time.Time `json:"last_password_change"`
	VerifiedEmail      bool       `json:"verified_email"`
	VerifiedPhone      bool       `json:"verified_phone"`
}

type LoginResponse struct {
	Token     string     `json:"token,omitempty"`
	BanExpire *time.Time `json:"ban_expire,omitempty"`
}

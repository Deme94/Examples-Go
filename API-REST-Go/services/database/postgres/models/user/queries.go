package user

import (
	"API-REST/services/database/postgres/models/role"
	"API-REST/services/database/postgres/predicates"
	"errors"
	"strings"
	"time"
)

func (m *Model) GetAll(p *predicates.Predicates) ([]*User, error) {
	query := m.Db.Table("users").Select(
		"id",
		"username",
		"email",
		"first_name",
		"last_name",
	)
	if p != nil {
		query = predicates.Apply(query, p)
	}
	res, err := query.Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("users not found")
	}

	var usrs []*User
	for _, r := range res {

		// Check if nil values
		firstName := ""
		fn := r["first_name"]
		if fn != nil {
			firstName = fn.(string)
		}
		lastName := ""
		ln := r["last_name"]
		if ln != nil {
			lastName = ln.(string)
		}

		u := User{
			ID:        int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
			Username:  r["username"].(string),
			Email:     r["email"].(string),
			FirstName: firstName,
			LastName:  lastName,
		}

		usrs = append(usrs, &u)
	}

	return usrs, nil
}
func (m *Model) Get(id int) (*User, error) {
	res, err := m.Db.Table("users").Select(
		"created_at",
		"updated_at",
		"deleted_at",

		"username",
		"email",

		"nick",
		"first_name",
		"last_name",
		"phone",
		"address",

		"last_login",
		"last_password_change",
		"verified_mail",
		"verified_phone",
		"ban_date",
		"ban_expire",

		"roles.id as role_id",
		"roles.name as role_name",
	).
		Where("users.id", "=", id).
		LeftJoin("users_roles", "users.id", "=", "users_roles.user_id").
		LeftJoin("roles", "users_roles.role_id", "=", "roles.id").
		Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}

	r := res[0]

	// Check if nil values
	firstName := ""
	fn := r["first_name"]
	if fn != nil {
		firstName = fn.(string)
	}
	lastName := ""
	ln := r["last_name"]
	if ln != nil {
		lastName = ln.(string)
	}
	phone := ""
	ph := r["phone"]
	if ph != nil {
		phone = ph.(string)
	}
	address := ""
	addr := r["address"]
	if addr != nil {
		address = addr.(string)
	}
	updatedAt := time.Time{}
	ua := r["updated_at"]
	if ua != nil {
		updatedAt = ua.(time.Time)
	}
	deletedAt := time.Time{}
	da := r["deleted_at"]
	if da != nil {
		deletedAt = da.(time.Time)
	}
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}
	banDate := time.Time{}
	bd := r["ban_date"]
	if bd != nil {
		banDate = bd.(time.Time)
	}
	banExpire := time.Time{}
	be := r["ban_expire"]
	if be != nil {
		banExpire = be.(time.Time)
	}

	// Relations
	var roles []*role.Role
	if res[0]["role_id"] != nil {
		for i, _ := range res {
			roles = append(roles, &role.Role{
				ID:   int(res[i]["role_id"].(int64)),
				Name: res[i]["role_name"].(string),
			})
		}
	}

	u := User{
		ID:                 id,
		CreatedAt:          r["created_at"].(time.Time),
		UpdatedAt:          updatedAt,
		DeletedAt:          deletedAt,
		Username:           r["username"].(string),
		Email:              r["email"].(string),
		Nick:               r["nick"].(string),
		FirstName:          firstName,
		LastName:           lastName,
		Phone:              phone,
		Address:            address,
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
		VerifiedMail:       r["verified_mail"].(bool),
		VerifiedPhone:      r["verified_phone"].(bool),
		BanDate:            banDate,
		BanExpire:          banExpire,
		Roles:              roles,
	}

	return &u, nil
}
func (m *Model) GetByEmailWithPassword(email string) (*User, error) {
	email = strings.ToLower(email)

	res, err := m.Db.Table("users").Select("id", "created_at", "username", "password", "last_login", "last_password_change").Where("email", "=", email).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Check if nil values
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}

	u := User{
		ID:                 int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
		Username:           r["username"].(string),
		Email:              email,
		Password:           r["password"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (m *Model) GetByUsernameWithPassword(username string) (*User, error) {
	username = strings.ToLower(username)

	res, err := m.Db.Table("users").Select("id", "email", "password", "created_at", "last_login", "last_password_change").Where("LOWER(username)", "=", username).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Atributes can be nil
	lastLogin := time.Time{}
	ll := r["last_login"]
	if ll != nil {
		lastLogin = ll.(time.Time)
	}

	u := User{
		ID:                 int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
		Username:           username,
		Email:              r["email"].(string),
		Password:           r["password"].(string),
		CreatedAt:          r["created_at"].(time.Time),
		LastLogin:          lastLogin,
		LastPasswordChange: r["last_password_change"].(time.Time),
	}

	return &u, nil
}
func (m *Model) GetPhoto(id int) (string, error) {
	res, err := m.Db.Table("users").Select("photo_name").Where("id", "=", id).Get()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}

	if res[0]["photo_name"] == nil {
		return "", errors.New("user has not photo")
	}

	return res[0]["photo_name"].(string), nil
}
func (m *Model) GetCV(id int) (string, error) {
	res, err := m.Db.Table("users").Select("cv_name").Where("id", "=", id).Get()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}

	if res[0]["cv_name"] == nil {
		return "", errors.New("user has not CV")
	}

	return res[0]["cv_name"].(string), nil
}
func (m *Model) Insert(u *User) error {
	nick := u.Nick
	if u.Nick == "" {
		nick = u.Username
	}

	err := m.Db.Table("users").Insert(map[string]interface{}{
		"username": u.Username,
		"email":    strings.ToLower(u.Email),
		"password": u.Password,
		"nick":     nick,
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *Model) Update(u *User) error {
	_, err := m.Db.Table("users").Where("id", "=", u.ID).Update(map[string]interface{}{
		"nick":       u.Nick,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"phone":      u.Phone,
		"address":    u.Address,
	})
	return err
}
func (m *Model) UpdatePassword(id int, password string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"password": password, "last_password_change": "NOW()"})
	return err
}
func (m *Model) UpdateRoles(id int, roleIDs ...int) error {
	if len(roleIDs) == 0 {
		return errors.New("no role id parameter was given")
	}

	return m.Db.InTransaction(func() (interface{}, error) {
		table := m.Db.Table("users_roles")
		// Remove all user_roles
		_, err := table.Where("user_id", "=", id).Delete()
		if err != nil {
			return nil, err
		}
		// Add new user_roles
		user_role := make(map[string]interface{})
		user_role["user_id"] = id
		for _, roleID := range roleIDs {
			user_role["role_id"] = roleID
			err = table.Insert(user_role)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}
func (m *Model) UpdatePhoto(id int, photoName string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"photo_name": photoName})
	return err
}
func (m *Model) UpdateCV(id int, cvName string) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"cv_name": cvName})
	return err
}
func (m *Model) Delete(id int) error {
	_, err := m.Db.Table("users").Where("id", "=", id).Update(map[string]interface{}{"deleted_at": "NOW()"})
	return err
}

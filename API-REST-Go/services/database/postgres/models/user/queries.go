package user

import (
	"API-REST/services/database/postgres/models/permission"
	"API-REST/services/database/postgres/models/role"
	"API-REST/services/database/postgres/predicates"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
			ID:        uuid.MustParse(r["id"].(string)), // la DB devuelve interface{} y se hace cast a int
			Username:  r["username"].(string),
			Email:     r["email"].(string),
			FirstName: firstName,
			LastName:  lastName,
		}

		usrs = append(usrs, &u)
	}

	return usrs, nil
}
func (m *Model) Get(id uuid.UUID) (*User, error) {
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
		"verified_email",
		"verified_phone",
		"ban_date",
		"ban_expire",

		"roles.id as role_id",
		"roles.name as role_name",
	).
		Where("users.id", "=", id.String()).
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
	var updatedAt *time.Time
	ua := r["updated_at"]
	if ua != nil {
		t := ua.(time.Time)
		updatedAt = &t
	}
	var deletedAt *time.Time
	da := r["deleted_at"]
	if da != nil {
		t := da.(time.Time)
		deletedAt = &t
	}
	var lastLogin *time.Time
	ll := r["last_login"]
	if ll != nil {
		t := ll.(time.Time)
		lastLogin = &t
	}
	var banDate *time.Time
	bd := r["ban_date"]
	if bd != nil {
		t := bd.(time.Time)
		banDate = &t
	}
	var banExpire *time.Time
	be := r["ban_expire"]
	if be != nil {
		t := be.(time.Time)
		banExpire = &t

	}

	// Relations
	var roles []*role.Role
	if res[0]["role_id"] != nil {
		for _, row := range res {
			roles = append(roles, &role.Role{
				ID:   int(row["role_id"].(int64)),
				Name: row["role_name"].(string),
			})
		}
	}

	createdAt := r["created_at"].(time.Time)
	LastPasswordChange := r["last_password_change"].(time.Time)

	u := User{
		ID:                 id,
		CreatedAt:          &createdAt,
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
		LastPasswordChange: &LastPasswordChange,
		VerifiedEmail:      r["verified_email"].(bool),
		VerifiedPhone:      r["verified_phone"].(bool),
		BanDate:            banDate,
		BanExpire:          banExpire,
		Roles:              roles,
	}

	return &u, nil
}
func (m *Model) GetIDByEmail(email string) (uuid.UUID, error) {
	res, err := m.Db.Table("users").Select(
		"id",
	).
		Where("email", "=", email).
		First()
	if err != nil {
		return uuid.Nil, err
	}
	if len(res) == 0 {
		return uuid.Nil, errors.New("user not found")
	}
	id := uuid.MustParse(res["id"].(string))

	return id, nil
}
func (m *Model) GetPassword(id uuid.UUID) (string, error) {
	res, err := m.Db.Table("users").Select(
		"password",
	).
		Where("id", "=", id.String()).
		First()
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("user not found")
	}
	password := res["password"].(string)

	return password, nil
}
func (m *Model) GetByEmailWithPassword(email string) (*User, error) {
	email = strings.ToLower(email)

	res, err := m.Db.Table("users").Select(
		"id",
		"created_at",
		"deleted_at",
		"username",
		"password",
		"last_login",
		"last_password_change",
		"ban_date",
		"ban_expire",
	).
		Where("email", "=", email).
		Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Check if nil values
	var lastLogin *time.Time
	ll := r["last_login"]
	if ll != nil {
		t := ll.(time.Time)
		lastLogin = &t
	}
	var deletedAt *time.Time
	da := r["deleted_at"]
	if da != nil {
		t := da.(time.Time)
		deletedAt = &t
	}
	var banDate *time.Time
	bd := r["ban_date"]
	if bd != nil {
		t := bd.(time.Time)
		banDate = &t
	}
	var banExpire *time.Time
	be := r["ban_expire"]
	if be != nil {
		t := be.(time.Time)
		banExpire = &t
	}

	createdAt := r["created_at"].(time.Time)
	lastPasswordChange := r["last_password_change"].(time.Time)

	u := User{
		ID:                 uuid.MustParse(r["id"].(string)), // la DB devuelve interface{} y se hace cast a int
		CreatedAt:          &createdAt,
		DeletedAt:          deletedAt,
		Username:           r["username"].(string),
		Email:              email,
		Password:           r["password"].(string),
		LastLogin:          lastLogin,
		LastPasswordChange: &lastPasswordChange,
		BanDate:            banDate,
		BanExpire:          banExpire,
	}

	return &u, nil
}
func (m *Model) GetByUsernameWithPassword(username string) (*User, error) {
	username = strings.ToLower(username)

	res, err := m.Db.Table("users").Select(
		"id",
		"created_at",
		"deleted_at",
		"email",
		"password",
		"last_login",
		"last_password_change",
		"ban_date",
		"ban_expire",
	).
		Where("LOWER(username)", "=", username).
		Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("user not found")
	}
	r := res[0]

	// Check if nil values
	var lastLogin *time.Time
	ll := r["last_login"]
	if ll != nil {
		t := ll.(time.Time)
		lastLogin = &t
	}
	var deletedAt *time.Time
	da := r["deleted_at"]
	if da != nil {
		t := da.(time.Time)
		deletedAt = &t
	}
	var banDate *time.Time
	bd := r["ban_date"]
	if bd != nil {
		t := bd.(time.Time)
		banDate = &t
	}
	var banExpire *time.Time
	be := r["ban_expire"]
	if be != nil {
		t := be.(time.Time)
		banExpire = &t
	}

	createdAt := r["created_at"].(time.Time)
	lastPasswordChange := r["last_password_change"].(time.Time)

	u := User{
		ID:                 uuid.MustParse(r["id"].(string)), // la DB devuelve interface{} y se hace cast a int
		CreatedAt:          &createdAt,
		DeletedAt:          deletedAt,
		Username:           username,
		Email:              r["email"].(string),
		Password:           r["password"].(string),
		LastLogin:          lastLogin,
		LastPasswordChange: &lastPasswordChange,
		BanDate:            banDate,
		BanExpire:          banExpire,
	}

	return &u, nil
}
func (m *Model) HasPermission(id uuid.UUID, p *permission.Permission) (bool, error) {
	res, err := m.Db.Table("users").Select(
		"permissions.id as permission_id",
		"permissions.resource as permission_resource",
		"permissions.operation as permission_operation",
	).
		Where("users.id", "=", id.String()).
		AndWhere("permissions.resource", "=", p.Resource).
		AndWhere("permissions.operation", "=", p.Operation).
		LeftJoin("users_roles", "users.id", "=", "users_roles.user_id").
		LeftJoin("roles", "users_roles.role_id", "=", "roles.id").
		LeftJoin("roles_permissions", "roles.id", "=", "roles_permissions.role_id").
		LeftJoin("permissions", "roles_permissions.permission_id", "=", "permissions.id").
		Get()
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}

	return true, nil
}
func (m *Model) HasVerifiedEmail(id uuid.UUID) (bool, error) {
	res, err := m.Db.Table("users").Select(
		"verified_email",
	).
		Where("users.id", "=", id.String()).
		Get()
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	verifiedEmail := res[0]["verified_email"].(bool)
	return verifiedEmail, nil
}
func (m *Model) HasVerifiedPhone(id uuid.UUID) (bool, error) {
	res, err := m.Db.Table("users").Select(
		"verified_phone",
	).
		Where("users.id", "=", id.String()).
		Get()
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	verifiedPhone := res[0]["verified_phone"].(bool)
	return verifiedPhone, nil
}
func (m *Model) GetPhoto(id uuid.UUID) (string, error) {
	res, err := m.Db.Table("users").Select("photo_name").Where("id", "=", id.String()).Get()
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
func (m *Model) GetCV(id uuid.UUID) (string, error) {
	res, err := m.Db.Table("users").Select("cv_name").Where("id", "=", id.String()).Get()
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

	// Begin transaction
	tx, err := m.Db.Sql().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert user
	columns := "id, username, email, password, nick"
	values := "'" + uuid.NewString() + "', '" + strings.ToLower(u.Username) + "', '" + strings.ToLower(u.Email) + "', '" + u.Password + "', '" + nick + "'"
	if u.FirstName != "" {
		columns += ", first_name"
		values += ", '" + u.FirstName + "'"
	}
	if u.LastName != "" {
		columns += ", last_name"
		values += ", '" + u.LastName + "'"
	}
	if u.Phone != "" {
		columns += ", phone"
		values += ", '" + u.Phone + "'"
	}
	if u.Address != "" {
		columns += ", address"
		values += ", '" + u.Address + "'"
	}
	_, err = tx.Exec("INSERT INTO users (" + columns + ") VALUES (" + values + ");")
	if err != nil {
		return err
	}

	// Check first user created
	var count int
	row := tx.QueryRow("SELECT COUNT(*) FROM users")
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 1 {
		// get last inserted user id
		var userID string
		row := tx.QueryRow("SELECT id FROM users LIMIT 1")
		err = row.Scan(&userID)
		if err != nil {
			return err
		}

		// get superadmin role id
		var superadminRoleID int
		row = tx.QueryRow("SELECT id FROM roles WHERE name = 'superadmin';")
		err = row.Scan(&superadminRoleID)
		if err != nil {
			return err
		}

		// Assign role to user (insert users_roles)
		_, err = tx.Exec("INSERT INTO users_roles (user_id, role_id) VALUES " +
			"('" + userID + "', " + fmt.Sprint(superadminRoleID) + ");")
		if err != nil {
			return err
		}
	}
	// Commit
	return tx.Commit()
}
func (m *Model) Update(u *User) error {
	colValues := make(map[string]interface{})
	if u.Nick != "" {
		colValues["nick"] = u.Nick
	}
	if u.FirstName != "" {
		colValues["first_name"] = u.FirstName
	}
	if u.LastName != "" {
		colValues["last_name"] = u.LastName
	}
	if u.Phone != "" {
		colValues["phone"] = u.Phone
	}
	if u.Address != "" {
		colValues["address"] = u.Address
	}
	if len(colValues) > 0 {
		colValues["updated_at"] = "NOW()"
	}
	_, err := m.Db.Table("users").Where("id", "=", u.ID.String()).Update(colValues)
	return err
}
func (m *Model) UpdatePassword(id uuid.UUID, password string) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"password": password, "last_password_change": "NOW()"})
	return err
}
func (m *Model) UpdateRoles(id uuid.UUID, roleIDs ...int) error {
	if len(roleIDs) == 0 {
		return errors.New("no role id parameter was given")
	}

	// Begin transaction
	tx, err := m.Db.Sql().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear user roles
	_, err = tx.Exec("DELETE from users_roles WHERE user_id = '" + id.String() + "';")
	if err != nil {
		return err
	}

	// Assign roles to user (insert users_roles)
	values := ""
	for _, roleID := range roleIDs {
		values += "('" + id.String() + "', " + fmt.Sprint(roleID) + "),"
	}
	values = strings.TrimSuffix(values, ",")

	_, err = tx.Exec("INSERT INTO users_roles (user_id, role_id) VALUES " + values + ";")
	if err != nil {
		return err
	}

	// Commit
	return tx.Commit()
}
func (m *Model) UpdatePhoto(id uuid.UUID, photoName string) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"photo_name": photoName})
	return err
}
func (m *Model) UpdateCV(id uuid.UUID, cvName string) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"cv_name": cvName})
	return err
}
func (m *Model) VerifyEmail(id uuid.UUID) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"verified_email": "true"})
	return err
}
func (m *Model) Ban(id uuid.UUID, banExpire time.Time) error {
	sqlStatement := `
		UPDATE users
		SET ban_date = $2, ban_expire = $3
		WHERE id = $1;`

	_, err := m.Db.Sql().Exec(sqlStatement, id.String(), "NOW()", banExpire)

	return err
}
func (m *Model) Unban(id uuid.UUID) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{
			"ban_date":   nil,
			"ban_expire": nil,
		})
	return err
}
func (m *Model) Restore(id uuid.UUID) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"deleted_at": nil})
	return err
}
func (m *Model) Delete(id uuid.UUID) error {
	_, err := m.Db.Table("users").
		Where("id", "=", id.String()).
		Update(map[string]interface{}{"deleted_at": "NOW()"})
	return err
}

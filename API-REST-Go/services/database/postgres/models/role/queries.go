package role

import (
	"API-REST/services/database/postgres/models/permission"
	"API-REST/services/database/postgres/predicates"
	"errors"
	"strings"
)

func (m *Model) GetAll(p *predicates.Predicates) ([]*Role, error) {
	query := m.Db.Table("roles").Select(
		"id",
		"name",
	)
	query = predicates.Apply(query, p)
	if p.HasWhere() {
		query = query.AndWhere("roles.name", "!=", "superadmin")
	} else {
		query = query.Where("roles.name", "!=", "superadmin")
	}

	res, err := query.Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("roles not found")
	}

	var roles []*Role
	for _, r := range res {
		p := Role{
			ID:   int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
			Name: r["name"].(string),
		}

		roles = append(roles, &p)
	}

	return roles, nil
}
func (m *Model) Get(id int) (*Role, error) {
	res, err := m.Db.Table("roles").Select(
		"name",

		"permissions.id as permission_id",
		"permissions.resource as permission_resource",
		"permissions.operation as permission_operation",
	).Where("roles.id", "=", id).AndWhere("roles.name", "!=", "superadmin").
		LeftJoin("roles_permissions", "roles.id", "=", "roles_permissions.role_id").
		LeftJoin("permissions", "roles_permissions.permission_id", "=", "permissions.id").
		Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("role not found")
	}

	r := res[0]

	// Relations
	var permissions []*permission.Permission
	if res[0]["permission_id"] != nil {
		for _, row := range res {
			permissions = append(permissions, &permission.Permission{
				ID:        int(row["permission_id"].(int64)),
				Resource:  row["permission_resource"].(string),
				Operation: row["permission_operation"].(string),
			})
		}
	}

	role := Role{
		ID:          id,
		Name:        r["name"].(string),
		Permissions: permissions,
	}

	return &role, nil
}
func (m *Model) Insert(r *Role) error {
	err := m.Db.Table("roles").Insert(map[string]interface{}{
		"name": strings.ToLower(r.Name),
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *Model) Update(r *Role) error {
	_, err := m.Db.Table("roles").Where("id", "=", r.ID).AndWhere("roles.name", "!=", "superadmin").
		Update(map[string]interface{}{
			"name": strings.ToLower(r.Name),
		})
	return err
}
func (m *Model) UpdatePermissions(id int, permissionIDs ...int) error {
	if len(permissionIDs) == 0 {
		return errors.New("no permission id parameter was given")
	}

	return m.Db.InTransaction(func() (interface{}, error) {
		table := m.Db.Table("roles_permissions")
		// Remove all role_permissions
		_, err := table.Where("role_id", "=", id).AndWhere("roles.name", "!=", "superadmin").Delete()
		if err != nil {
			return nil, err
		}
		// Add new role_permissions
		role_permission := make(map[string]interface{})
		role_permission["role_id"] = id
		for _, permissionID := range permissionIDs {
			role_permission["permission_id"] = permissionID
			err = table.Insert(role_permission)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}
func (m *Model) Delete(id int) error {
	_, err := m.Db.Table("roles").Where("id", "=", id).AndWhere("roles.name", "!=", "superadmin").
		Delete()
	return err
}

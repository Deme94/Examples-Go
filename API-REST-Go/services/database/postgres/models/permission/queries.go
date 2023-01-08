package permission

import (
	"API-REST/services/database/postgres/predicates"
	"errors"
	"strings"
)

func (m *Model) GetAll(p *predicates.Predicates) ([]*Permission, error) {
	query := m.Db.Table("permissions").Select(
		"id",
		"resource",
		"operation",
	)
	if p != nil {
		query = predicates.Apply(query, p)
	}
	res, err := query.Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("permissions not found")
	}

	var perms []*Permission
	for _, r := range res {
		p := Permission{
			ID:        int(r["id"].(int64)), // la DB devuelve interface{} y se hace cast a int
			Resource:  r["resource"].(string),
			Operation: r["operation"].(string),
		}

		perms = append(perms, &p)
	}

	return perms, nil
}
func (m *Model) Get(id int) (*Permission, error) {
	res, err := m.Db.Table("permissions").Select(
		"resource",
		"operation",
	).Where("id", "=", id).Get()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("permission not found")
	}

	r := res[0]

	p := Permission{
		ID:        id,
		Resource:  r["resource"].(string),
		Operation: r["operation"].(string),
	}

	return &p, nil
}
func (m *Model) Insert(p *Permission) error {
	err := m.Db.Table("permissions").Insert(map[string]interface{}{
		"resource":  strings.ToLower(p.Resource),
		"operation": strings.ToLower(p.Operation),
	})
	if err != nil {
		return err
	}
	return nil
}
func (m *Model) Update(p *Permission) error {
	_, err := m.Db.Table("permissions").Where("id", "=", p.ID).Update(map[string]interface{}{
		"resource":  strings.ToLower(p.Resource),
		"operation": strings.ToLower(p.Operation),
	})
	return err
}
func (m *Model) Delete(id int) error {
	_, err := m.Db.Table("permissions").Where("id", "=", id).Delete()
	return err
}

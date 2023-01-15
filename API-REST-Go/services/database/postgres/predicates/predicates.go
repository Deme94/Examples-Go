package predicates

import (
	"strings"

	"github.com/arthurkushman/buildsqlx"
)

type where struct {
	and        bool
	or         bool
	columnName string
	operator   string
	value      interface{}
}

type orderBy struct {
	columnName string
	desc       bool
}

type Predicates struct {
	wheres  []*where
	groupBy string
	orderBy *orderBy
	limit   int
	offset  int
}

// Where
func (p *Predicates) Where(columnName string, operator string, value interface{}) *Predicates {
	p.wheres = append(p.wheres, &where{columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) AndWhere(columnName string, operator string, value interface{}) *Predicates {
	p.wheres = append(p.wheres, &where{and: true, columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) OrWhere(columnName string, operator string, value interface{}) *Predicates {
	p.wheres = append(p.wheres, &where{or: true, columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) WhereInsensitive(columnName string, operator string, value interface{}) *Predicates {
	columnName = "LOWER(" + columnName + ")"
	value = strings.ToLower(value.(string))
	p.wheres = append(p.wheres, &where{columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) AndWhereInsensitive(columnName string, operator string, value interface{}) *Predicates {
	columnName = "LOWER(" + columnName + ")"
	value = strings.ToLower(value.(string))
	p.wheres = append(p.wheres, &where{and: true, columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) OrWhereInsensitive(columnName string, operator string, value interface{}) *Predicates {
	columnName = "LOWER(" + columnName + ")"
	value = strings.ToLower(value.(string))
	p.wheres = append(p.wheres, &where{or: true, columnName: columnName, operator: operator, value: value})
	return p
}
func (p *Predicates) HasWhere() bool {
	return len(p.wheres) > 0
}

// GroupBy
func (p *Predicates) GroupBy(columnName string) *Predicates {
	p.groupBy = columnName
	return p
}

// OrderBy
func (p *Predicates) OrderBy(columnName string, desc bool) *Predicates {
	p.orderBy = &orderBy{columnName: columnName, desc: desc}
	return p
}

// Limit
func (p *Predicates) Limit(limit int) *Predicates {
	p.limit = limit
	return p
}

// Offset
func (p *Predicates) Offset(offset int) *Predicates {
	p.offset = offset
	return p
}

func Apply(query *buildsqlx.DB, p *Predicates) *buildsqlx.DB {
	for _, where := range p.wheres {
		if where.and {
			query = query.AndWhere(where.columnName, where.operator, where.value)
		} else if where.or {
			query = query.OrWhere(where.columnName, where.operator, where.value)
		} else {
			query = query.Where(where.columnName, where.operator, where.value)
		}
	}
	if p.groupBy != "" {
		query = query.GroupBy(p.groupBy)
	}
	if p.orderBy != nil {
		dir := "ASC"
		if p.orderBy.desc {
			dir = "DESC"
		}
		query = query.OrderBy(p.orderBy.columnName, dir)
	}
	if p.limit != 0 {
		query = query.Limit(int64(p.limit))
	}
	if p.offset != 0 {
		query = query.Offset(int64(p.offset))
	}

	return query
}

package gokit

import "gorm.io/gorm"

type GormQueryBuilder struct {
	db    *gorm.DB
	model any
}

func NewGormQueryBuilder(db *gorm.DB, model any) QueryBuilder {
	return &GormQueryBuilder{
		db:    db,
		model: model,
	}
}

func (g *GormQueryBuilder) Select(columns ...string) QueryBuilder {
	g.db = g.db.Select(columns)
	return g
}

func (g *GormQueryBuilder) Where(query any, args ...any) QueryBuilder {
	g.db = g.db.Where(query, args...)
	return g
}

func (g *GormQueryBuilder) WhereIn(column string, values []any) QueryBuilder {
	g.db = g.db.Where(column+" IN ?", values)
	return g
}

func (g *GormQueryBuilder) Join(query string, args ...any) QueryBuilder {
	g.db = g.db.Joins(query, args...)
	return g
}

func (g *GormQueryBuilder) OrderBy(column string, direction ...string) QueryBuilder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = direction[0]
	}
	g.db = g.db.Order(column + " " + dir)
	return g
}

func (g *GormQueryBuilder) GroupBy(columns ...string) QueryBuilder {
	for _, col := range columns {
		g.db = g.db.Group(col)
	}
	return g
}

func (g *GormQueryBuilder) Having(query any, args ...any) QueryBuilder {
	g.db = g.db.Having(query, args...)
	return g
}

func (g *GormQueryBuilder) Limit(limit int) QueryBuilder {
	g.db = g.db.Limit(limit)
	return g
}

func (g *GormQueryBuilder) Offset(offset int) QueryBuilder {
	g.db = g.db.Offset(offset)
	return g
}

func (g *GormQueryBuilder) Preload(query string, args ...any) QueryBuilder {
	g.db = g.db.Preload(query, args...)
	return g
}

func (g *GormQueryBuilder) Joins(query string, args ...any) QueryBuilder {
	g.db = g.db.Joins(query, args...)
	return g
}

func (g *GormQueryBuilder) Find(dest any) error {
	return g.db.Find(dest).Error
}

func (g *GormQueryBuilder) First(dest any) error {
	return g.db.First(dest).Error
}

func (g *GormQueryBuilder) Create(value any) error {
	return g.db.Create(value).Error
}

func (g *GormQueryBuilder) Update(values any) error {
	return g.db.Updates(values).Error
}

func (g *GormQueryBuilder) Delete() error {
	return g.db.Delete(g.model).Error
}

func (g *GormQueryBuilder) Count() (int64, error) {
	var count int64
	err := g.db.Model(g.model).Count(&count).Error
	return count, err
}

func (g *GormQueryBuilder) Begin() (QueryBuilder, error) {
	tx := g.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return NewGormQueryBuilder(tx, g.model), nil
}

func (g *GormQueryBuilder) Commit() error {
	return g.db.Commit().Error
}

func (g *GormQueryBuilder) Rollback() error {
	return g.db.Rollback().Error
}

func (g *GormQueryBuilder) ToSQL() (string, []any, error) {
	stmt := &gorm.Statement{DB: g.db}
	g.db.Statement = stmt

	return g.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(g.model)
	}), nil, nil
}

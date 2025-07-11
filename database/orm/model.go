package orm

import "time"

type Model struct {
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" db:"updated_at"`
}

type SoftDeletes struct {
	Model
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

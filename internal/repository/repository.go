package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) CountByField(db *gorm.DB, field string, value interface{}) (int64, error) {
	var count int64
	err := db.Unscoped().Model(new(T)).Where(field+" = ?", value).Count(&count).Error
	return count, err
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, entity *T, email string) error {
	return db.Where("email = ?", email).Take(entity).Error
}
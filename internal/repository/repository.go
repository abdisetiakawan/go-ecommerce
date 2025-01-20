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

func (r *Repository[T]) HasUserID(db *gorm.DB, userID uint) (bool, error) {
	var count int64
	err := db.Model(new(T)).Where("user_id = ?", userID).Count(&count).Error
	return count > 0, err
}

func (r *Repository[T]) FindByID(db *gorm.DB, entity *T, id uint) error {
	return db.First(entity, id).Error
}

func (r *Repository[T]) FindByUUID(db *gorm.DB, entity *T, uuid string, fieldName string) error {
	return db.Where(fieldName+" = ?", uuid).First(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, id uint) error {
    return db.Unscoped().Delete(new(T), id).Error
}

func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T, limit, offset int) error {
	return db.Limit(limit).Offset(offset).Find(entities).Error
}

func (r *Repository[T]) FindByStatus(db *gorm.DB, entities *[]T, status string) error {
	return db.Where("status = ?", status).Find(entities).Error
}

func (r *Repository[T]) UpdateStatus(db *gorm.DB, id uint, status string) error {
	return db.Model(new(T)).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository[T]) FindByUserID(db *gorm.DB, entities *[]T, userID uint) error {
	return db.Where("user_id = ?", userID).Find(entities).Error
}

func (r *Repository[T]) FindByStoreID(db *gorm.DB, entities *[]T, storeID uint) error {
	return db.Where("store_id = ?", storeID).Find(entities).Error
}

func (r *Repository[T]) CountAll(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(new(T)).Count(&count).Error
	return count, err
}

func (r *Repository[T]) ExistsByField(db *gorm.DB, field string, value interface{}) (bool, error) {
	var exists bool
	err := db.Model(new(T)).Select("count(*) > 0").Where(field+" = ?", value).Find(&exists).Error
	return exists, err
}

func (r *Repository[T]) SoftDelete(db *gorm.DB, id uint) error {
	return db.Delete(new(T), id).Error
}

func (r *Repository[T]) Restore(db *gorm.DB, id uint) error {
	return db.Unscoped().Model(new(T)).Where("id = ?", id).Update("deleted_at", nil).Error
}

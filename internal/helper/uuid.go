package helper

import "github.com/google/uuid"


type UUIDHelper struct {
	Value string
}
func NewUUIDHelper() *UUIDHelper {
	return &UUIDHelper{Value: uuid.New().String()}
}
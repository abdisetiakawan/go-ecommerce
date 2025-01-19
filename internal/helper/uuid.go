package helper

import "github.com/google/uuid"

type UUIDHelper struct{}

func NewUUIDHelper() *UUIDHelper {
    return &UUIDHelper{}
}

func (u *UUIDHelper) Generate() string {
    return uuid.New().String()
}

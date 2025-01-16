package model

import "net/http"

type WebResponse[T any] struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    T             `json:"data"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Errors  string        `json:"errors,omitempty"`
}

func NewWebResponse[T any](data T, message string, statusCode int, paging *PageMetadata) WebResponse[T] {
	status := "fail"
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		status = "success"
	}

	return WebResponse[T]{
		Status:  status,
		Message: message,
		Data:    data,
		Paging:  paging,
	}
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetadata PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

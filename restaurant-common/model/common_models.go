package model

type CommonErrorType string

const (
	CommonErrorTypeWrongRequest CommonErrorType = "Wrong request"
	CommonErrorTypeInvalidData  CommonErrorType = "Invalid data"
)

type CommonErrorResponse struct {
	Type    CommonErrorType `json:"type"`
	Message string          `json:"message"`
}

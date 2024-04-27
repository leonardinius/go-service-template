package apiworker

import (
	sharedv1 "github.com/leonardinius/go-service-template/internal/apigen/shared/v1"
)

func errToSharedError(err error) *sharedv1.Error {
	if err == nil {
		return nil
	}

	return &sharedv1.Error{
		Code:    500,
		Message: err.Error(),
	}
}

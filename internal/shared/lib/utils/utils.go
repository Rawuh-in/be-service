package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	paginationModel "rawuh-service/internal/shared/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func Base64DecodeStripped(s string) (string, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	fmt.Println(encoded)

	return string(decoded), err
}

func IsValidCharacter(input string) bool {
	regex := `^[\s\w\d_.,-;()/]*$`
	match, err := regexp.MatchString(regex, input)
	if err != nil {
		return false
	}
	return match
}

func IsValidProductName(input string) bool {
	regex := `^[a-zA-Z0-9 _.,'-]+$`
	match, err := regexp.MatchString(regex, input)
	if err != nil {
		return false
	}
	return match

}

func IsEmptyString(value string) bool {
	return strings.TrimSpace(value) == ""
}

type APIErrorResponse struct {
	Error   bool   `json:"Error"`
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}

// HandleGrpcError converts a gRPC error into a proper HTTP JSON response
func HandleGrpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if ok {
		// Map gRPC codes to HTTP status
		var httpCode int
		switch st.Code() {
		case codes.InvalidArgument:
			httpCode = http.StatusBadRequest
		case codes.NotFound:
			httpCode = http.StatusNotFound
		case codes.PermissionDenied:
			httpCode = http.StatusForbidden
		case codes.Unauthenticated:
			httpCode = http.StatusUnauthorized
		case codes.AlreadyExists:
			httpCode = http.StatusConflict
		default:
			httpCode = http.StatusInternalServerError
		}

		resp := APIErrorResponse{
			Error:   true,
			Code:    httpCode,
			Message: st.Message(), // clean message only
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpCode)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Non-gRPC error fallback
	resp := APIErrorResponse{
		Error:   true,
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(resp)
}

func SetPagination(page int32, limit int32) *paginationModel.PaginationResponse {
	res := &paginationModel.PaginationResponse{
		Limit: 10,
		Page:  1,
	}

	if limit == 0 && page == 0 {
		res.Limit = -1
		res.Page = -1
		return res
	} else {
		res.Limit = limit
		res.Page = page
	}

	if res.Page == 0 {
		res.Page = 1
	}

	switch {
	case res.Limit > 100:
		res.Limit = 100
	case res.Limit <= 0:
		res.Limit = 10
	}

	return res
}

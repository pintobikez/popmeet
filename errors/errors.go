package errors

type ErrResponse struct {
	Error ErrContent `json:"error"`
}

type ErrContent struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HealthStatus struct {
	Repo *HealthStatusDetail `json:"repository"`
}

type HealthStatusDetail struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

const (
	ErrorInterestNotFound    = 1001
	ErrorInterestsNotFound   = 1002
	ErrorUserNotFound        = 1003
	ErrorUserProfileNotFound = 1004
	ErrorCreatingToken       = 1005
	ErrorEventNotFound       = 1006
	ErrorCantAddUSerToEvent  = 1007
)

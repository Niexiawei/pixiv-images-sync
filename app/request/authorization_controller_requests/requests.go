package authorization_controller_requests

type Receive struct {
	Code         string `json:"code" form:"code" binding:"required"`
	State        string `json:"state" form:"state" binding:"omitempty"`
	SessionState string `json:"session_state" form:"session_state" binding:"required"`
}

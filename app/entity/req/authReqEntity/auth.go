package authReqEntity

// /api/auth/
type AuthInput struct {
	//Token        string
	Token string `form:"token"  binding:"required" msg:"必填"`
	//Action string `form:"action"  binding:"required" msg:"必填"`
}

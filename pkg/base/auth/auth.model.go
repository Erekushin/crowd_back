package auth

type LoginReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Serialno    string `json:"serial_no"`
	DeviceToken string `json:"device_token"`
}

type RegisterReq struct {
	Identity string `json:"identity" validate:"required"`
	Password string `json:"password" validate:"required"`
	Otp      uint   `json:"otp,string" validate:"required"`
}

package account_model

type RespCgToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

type ReqCgWithdraw struct {
	SrcAcc      string  `json:"srcAcc"`
	DestAcc     string  `json:"destAcc"`
	DestName    string  `json:"destName"`
	DestBank    string  `json:"destBank"`
	Amount      float32 `json:"amount"`
	Description string  `json:"description"`
}

package account_model

type ReqGolomtConfirmation struct {
	Callback      string `json:"callback"`
	Checksum      string `json:"checksum"`
	ReturnType    string `json:"returnType"`
	TransactionId string `json:"transactionId"`
}

type RespGolomtConfirmation struct {
	Checksum      string `json:"checksum"`
	Invoice       string `json:"invoice"`
	TransactionId string `json:"transactionId"`
}

type ReqGolomtPayment struct {
	Amount        float32 `json:"amount"`
	Checksum      string  `json:"checksum"`
	TransactionId string  `json:"transactionId"`
	Lang          string  `json:"lang"`
	Token         string  `json:"token"`
}

type RespGolomtPayment struct {
	Amount        string `json:"amount"`
	ErrorDesc     string `json:"errorDesc"`
	Checksum      string `json:"checksum"`
	ErrorCode     string `json:"errorCode"`
	TransactionId string `json:"transactionId"`
	CardNumber    string `json:"cardNumber"`
}

type ReqInvoice struct {
	Invoice     string `json:"invoice"`
	RedirectUrl string `json:"redirect_url"`
}

type ReqGolomtGetToken struct {
	Checksum      string `json:"checksum"`
	TransactionId string `json:"transactionId"`
}

type RespGolomtGetToken struct {
	BankCode      string `json:"bankCode"`
	Bank          string `json:"bank"`
	ErrorDesc     string `json:"errorDesc"`
	CheckSum      string `json:"checksum"`
	ErrorCode     string `json:"errorCode"`
	CardHolder    string `json:"cardHolder"`
	TransactionId string `json:"transactionId"`
	CardNumber    string `json:"cardNumber"`
	Token         string `json:"token"`
}

type StatementInfoReq struct {
	JournalNo string `json:"journal_no"`
	TranType  string `json:"tran_type"`
}

type ResBankTransaction struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  BankTxnResult `json:"result"`
}

type BankTxnResult struct {
	JournalNo string `json:"journalNo"`
	LogId     int    `json:"logId"`
}

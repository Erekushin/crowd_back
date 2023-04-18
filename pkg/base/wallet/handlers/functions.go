package account_handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	account_model "crowdfund/pkg/base/wallet/models"
	"crowdfund/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AccountDeposit(accountNo, refNo string, amount float32) error {
	userAccount, err := account_model.GetAccount(accountNo)
	if err != nil {
		return err
	}

	tdbDepositAccount := "2000000001"
	tdbMirrorAccount := "2000000003"

	journalNo := uuid.New().String()

	db := database.DBconn

	var txnsignLog account_model.TxnSignLog
	txnsignLog.Type = "+"
	txnsignLog.Amount = amount
	txnsignLog.BankJrNo = refNo
	txnsignLog.AppJrNo = journalNo
	txnsignLog.UserId = userAccount.OwnerId
	db.Save(&txnsignLog)

	if err := MakeTransaction(amount, journalNo, tdbDepositAccount, tdbMirrorAccount, "Цэнэглэлт", "", "", "", "", refNo); err != nil {
		return fmt.Errorf("transaction error1:" + err.Error())
	}

	if err := MakeTransaction(amount, journalNo, tdbMirrorAccount, userAccount.AccountNo, "Цэнэглэлт", "", "", "", "", refNo); err != nil {
		return fmt.Errorf("transaction error2:" + err.Error())
	}

	// notif_text := "Таны данс " + strconv.Itoa(int(amount)) + " төгрөгөөр цэнэглэгдлээ."
	// go client.SendNotif(userAccount.OwnerId, notif_text, notif_text, "txn", "")

	return nil
}

func depo_account(userId uint, amount float32, refNo, description string) error {

	fmt.Println("depo_account:", userId, amount, refNo, description)

	userAccount, err := account_model.GetDefaultAccount(userId)
	if err != nil {
		return err
	}

	golomtDepositAccount := "2000000002"
	golomtMirrorAccount := "2000000004"
	journalNo := uuid.New().String()

	db := database.DBconn

	var txnsignLog account_model.TxnSignLog
	txnsignLog.Type = "+"
	txnsignLog.Amount = amount
	txnsignLog.BankJrNo = refNo
	txnsignLog.AppJrNo = journalNo
	txnsignLog.UserId = userId

	if err := db.Save(&txnsignLog).Error; err != nil {
		fmt.Println("txnsignLog save:", err.Error())
	}

	if err := MakeTransaction(amount, journalNo, golomtDepositAccount, golomtMirrorAccount, description, "GOLOMT_CARD", "150000", "Голомт банк", description, refNo); err != nil {
		return fmt.Errorf("transaction error:" + err.Error())
	}

	if err := MakeTransaction(amount, journalNo, golomtMirrorAccount, userAccount.AccountNo, description, "GOLOMT_CARD", "150000", "Голомт банк", description, refNo); err != nil {
		return fmt.Errorf("transaction error")
	}

	// notif_text := "Таны данс " + strconv.Itoa(int(amount)) + " төгрөгөөр цэнэглэгдлээ."
	// go client.SendNotif(user_account.OwnerId, notif_text, notif_text, "txn", "")
	return err
}

func bankTransfer(req account_model.ReqCgWithdraw) (err error, ref_no string) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	r, err := http.NewRequest("POST", os.Getenv("CG_OAUTH_URL"), strings.NewReader(form.Encode()))
	r.Header.Set("Authorization", os.Getenv("CG_OAUTH_AUTH"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return err, ref_no
	}
	defer resp.Body.Close()

	fmt.Println("resp.Body:", resp.Body)

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ref_no
	}

	var token account_model.RespCgToken

	err = json.Unmarshal(result, &token)
	if err != nil {
		return err, ref_no
	}
	fmt.Println("\n\n\ntoken:", token)

	reqb, _ := json.Marshal(req)
	url := os.Getenv("CG_INTERBANK_TRANSFER_URL")
	if req.DestBank == "040000" {
		url = os.Getenv("CG_TDB_TRANSFER_URL")
	}
	req_cg, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqb))
	req_cg.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req_cg.Header.Set("Content-type", "application/json")
	resp, err = client.Do(req_cg)
	if err != nil {
		return err, ref_no
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ref_no
	}
	var bankResp account_model.ResBankTransaction
	err = json.Unmarshal([]byte(body), &bankResp)
	fmt.Println("Bank transaction bankResp:", bankResp)

	if bankResp.Status == "success" {
		return err, bankResp.Result.JournalNo
	} else {
		return fmt.Errorf("error"), ref_no
	}
}

func MakeTransaction(amount float32, journalNo, srcAccountNo, destAccountNo, description, paymentMethod, bankCode, bankName, bankAccount, refNo string) error {

	if paymentMethod == "" {
		paymentMethod = "wallet"
	}

	if bankCode == "" {
		bankCode = "000001"
	}

	if bankName == "" {
		bankName = "Gerege wallet"
	}

	var src_account = account_model.Account{}
	var dest_account = account_model.Account{}

	db := database.DBconn

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("account_no = ?", srcAccountNo).First(&src_account).Error
		if err != nil {
			return err
		}
		err = tx.Where("account_no = ?", destAccountNo).First(&dest_account).Error
		if err != nil {
			return err
		}
		var src_running_balance = src_account.Balance
		var dest_running_balance = dest_account.Balance

		dest_account.Balance = dest_account.Balance + amount
		src_account.Balance = src_account.Balance - amount

		tx.Save(&src_account)
		tx.Save(&dest_account)

		var txn = account_model.Transactions{
			JournalNo:          journalNo,
			SrcAccountNo:       srcAccountNo,
			SrcRunningBalance:  src_running_balance,
			DestAccountNo:      destAccountNo,
			DestRunningBalance: dest_running_balance,
			Description:        description,
			TranType:           "C",
			Amount:             amount,
			PaymentMethod:      paymentMethod,
			BankCode:           bankCode,
			BankName:           bankName,
			BankAccount:        bankAccount,
			RefNo:              refNo,
		}

		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		txn = account_model.Transactions{
			JournalNo:          journalNo,
			SrcAccountNo:       destAccountNo,
			SrcRunningBalance:  dest_running_balance,
			DestAccountNo:      srcAccountNo,
			DestRunningBalance: src_running_balance,
			Description:        description,
			TranType:           "D",
			Amount:             amount,
			PaymentMethod:      paymentMethod,
			BankCode:           bankCode,
			BankName:           bankName,
			BankAccount:        bankAccount,
			RefNo:              refNo,
		}

		if err := tx.Create(&txn).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func PasswordEncrypt(password string) string {
	sha1pwd := sha1.Sum([]byte(strings.ToUpper(password)))
	return strings.ToUpper(hex.EncodeToString(sha1pwd[:]))
}

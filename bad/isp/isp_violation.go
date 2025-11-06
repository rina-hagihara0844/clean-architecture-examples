// bad/payment_bad.go
package bad

import (
	"errors"
	"fmt"
	"time"
)

// 取引情報の最小モデル
type Transaction struct {
	ID       string
	Amount   int64
	Currency string
	Status   string
	Created  time.Time
}

// 振込情報の最小モデル
type Payout struct {
	ID         string
	MerchantID string
	Amount     int64
	Created    time.Time
}

type PaymentGateway interface {
	PayoutToMerchant(merchantID string, amount int64) error                              // 加盟店へ指定金額を振り込む処理
	GetTransaction(txID string) (Transaction, error)                                     // 取引IDから取引詳細を取得する処理
	ListPayouts(merchantID string, since time.Time) ([]Payout, error)                    // 加盟店の振込履歴を取得する処理
	GenerateInvoicePDF(invoiceID string) ([]byte, error)                                 // 請求書PDFを生成する処理
	SendReceiptEmail(toEmail, txID string) error                                         // 取引完了メールを送信する処理
	ChargeCard(cardToken string, amount int64, currency string) (txID string, err error) // カードで課金する処理
	Refund(txID string, amount int64) error                                              // 取引を返金する処理
	VerifyBankAccount(bankAccountID string) error                                        // 銀行口座の正当性を検証する処理（銀行系特有）
}

// 実装1：Stripe（カード決済メインだが、共通処理は一通り対応）
type StripeGateway struct{}

// 加盟店へ指定金額を振り込む処理（Stripe Connect 相当のつもりで成功扱い）
func (s StripeGateway) PayoutToMerchant(merchantID string, amount int64) error { return nil }

// 取引IDから取引詳細を取得する処理
func (s StripeGateway) GetTransaction(txID string) (Transaction, error) {
	return Transaction{ID: txID, Amount: 1000, Currency: "JPY", Status: "succeeded", Created: time.Now()}, nil
}

// 加盟店の振込履歴を取得する処理
func (s StripeGateway) ListPayouts(merchantID string, since time.Time) ([]Payout, error) {
	return []Payout{{ID: "po_1", MerchantID: merchantID, Amount: 5000, Created: time.Now()}}, nil
}

// 請求書PDFを生成する処理（ダミーPDF）
func (s StripeGateway) GenerateInvoicePDF(invoiceID string) ([]byte, error) {
	return []byte("%PDF-1.4 ...stripe..."), nil
}

// 取引完了メールを送信する処理
func (s StripeGateway) SendReceiptEmail(toEmail, txID string) error { return nil }

// カードで課金する処理
func (s StripeGateway) ChargeCard(cardToken string, amount int64, currency string) (string, error) {
	return fmt.Sprintf("tx_stripe_%d", time.Now().UnixNano()), nil
}

// 取引を返金する処理
func (s StripeGateway) Refund(txID string, amount int64) error { return nil }

// 銀行口座の正当性を検証する処理
func (s StripeGateway) VerifyBankAccount(bankAccountID string) error {
	return errors.New("not supported by StripeGateway: VerifyBankAccount")
}

// 実装2：BankTransfer（銀行振込メインだが、共通処理は一通り対応）
type BankTransferGateway struct{}

// 加盟店へ指定金額を振り込む処理
func (b BankTransferGateway) PayoutToMerchant(merchantID string, amount int64) error { return nil }

// 取引IDから取引詳細を取得する処理
func (b BankTransferGateway) GetTransaction(txID string) (Transaction, error) {
	return Transaction{ID: txID, Amount: 1000, Currency: "JPY", Status: "settled", Created: time.Now()}, nil
}

// 加盟店の振込履歴を取得する処理
func (b BankTransferGateway) ListPayouts(merchantID string, since time.Time) ([]Payout, error) {
	return []Payout{{ID: "po_b1", MerchantID: merchantID, Amount: 7000, Created: time.Now()}}, nil
}

// 請求書PDFを生成する処理（ダミーPDF）
func (b BankTransferGateway) GenerateInvoicePDF(invoiceID string) ([]byte, error) {
	return []byte("%PDF-1.4 ...bank..."), nil
}

// 取引完了メールを送信する処理（ダミー成功）
func (b BankTransferGateway) SendReceiptEmail(toEmail, txID string) error { return nil }

// カードで課金する処理
func (b BankTransferGateway) ChargeCard(cardToken string, amount int64, currency string) (string, error) {
	return "", errors.New("not supported by BankTransferGateway: ChargeCard")
}

// 取引を返金する処理
func (b BankTransferGateway) Refund(txID string, amount int64) error {
	return errors.New("not supported by BankTransferGateway: Refund")
}

// 銀行口座の正当性を検証する処理（得意）
func (b BankTransferGateway) VerifyBankAccount(bankAccountID string) error { return nil }

package mail

import "fmt"

// メール送信（ダミー実装）
func SendSMTP(to, body string) error {
	fmt.Println("to:", to)
	fmt.Println("body:", body)
	return nil
}

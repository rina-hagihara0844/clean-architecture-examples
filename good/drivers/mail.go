package drivers

// Framework & Drivers層（インフラストラクチャ）
// --------------------------------------------------------
// この層の責務：
// - 外部システム（mailer）との接続や入出力を扱う
// - Domain / UseCase 層のインターフェースを具象化して実装する
// -------------------------------------------------------
// - 業務ロジックを含まない（技術的な処理のみ）
// --------------------------------------------------------

type SMTPMailer struct{}

// メール送信の具象実装
func (m SMTPMailer) NotifyManagerNewRequest(id string) error { /* 実送信 */ return nil }

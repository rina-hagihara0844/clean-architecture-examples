# clean-architecture-examples

このリポジトリは、SOLID 原則を題材にした勉強会用のコード例をまとめたものです。良い設計へリファクタリングする前段として、あえてアンチパターンを残した状態のコードを題材にディスカッションできる構成になっています。

## 使っている言語

- Go（Go 1.21 以降を想定。標準ライブラリのみ使用）

## ローカルでの利用手順

1. リポジトリを取得します。
   ```bash
   git clone https://github.com/rina-hagihara0844/clean-architecture-examples.git
   cd clean-architecture-examples
   ```
2. Go がインストールされていることを確認します。
   ```bash
   go version
   ```
3. 各セクションのコードをエディターで確認しながら勉強会を進めます。`bad/server_bad.go` は `go run ./bad/server_bad.go` で実行可能ですが、外部サービス（DB・メール）がスタブ化されていないため、実行するとランタイムエラーが発生する想定です。動作確認ではなく、問題点の洗い出し用サンプルとして扱ってください。

## ファイルの説明

- `bad/server_bad.go` → オープン／クローズドの原則・依存性逆転の原則に違反しているサンプル。
- `bad/isp/isp_violation.go` → インターフェース分離の原則に違反しているサンプル。
- `good/` 配下 → リファクタリング後のクリーンアーキテクチャ例がまとまっています。
  勉強会後に答え合わせとして参照する用なので、途中では見ないでおくことを推奨します。

## 勉強会で扱う原則

- ⛔ 単一責任の原則（SRP）: 今回のコード例では扱いません。
- ✅ オープン／クローズドの原則（OCP）: `bad/server_bad.go`
- ✅ 依存性逆転の原則（DIP）: `bad/server_bad.go`
- ✅ インターフェース分離の原則（ISP）: `bad/isp/isp_violation.go`
- ⛔ リスコフの置換原則（LSP）: 今回のコード例では扱いません。

# go-load

CPU・メモリに対して指定した割合の負荷を一定時間かけるライブラリおよびCLIツール。
負荷試験・スケールテスト・監視アラートの動作確認などに使用できる。

## 機能

- CPU 使用率を指定割合に引き上げる（コア数分の並列ゴルーチン）
- メモリ使用量を指定割合まで増加させる
- 指定した秒数・分数・時間後に自動終了
- ライブラリとして Go コードから呼び出し可能

## インストール

### CLI ツール

```bash
go install github.com/rssh-jp/go-load/cmd/load@latest
```

### ライブラリ

```bash
go get github.com/rssh-jp/go-load
```

## 使い方

### CLI

```bash
# CPU 70%、メモリ 60%、30秒間負荷をかける
load -cpu 70 -mem 60 -duration 30s

# 1分間
load -cpu 50 -mem 50 -duration 1m

# オプション一覧
load -h
```

| フラグ | デフォルト | 説明 |
|--------|-----------|------|
| `-cpu` | `50` | CPU 使用率 (%) |
| `-mem` | `50` | メモリ使用率 (%) |
| `-duration` | `60s` | 負荷をかける時間 (`30s`, `1m`, `2h` 形式) |

### ライブラリ

```go
import (
    "time"
    load "github.com/rssh-jp/go-load"
)

l := load.New(
    load.OptionCPUPercentage(70),
    load.OptionMemoryPercentage(60),
    load.OptionDuration(30 * time.Second),
)

if err := l.Run(); err != nil {
    log.Fatal(err)
}
```

## 要件

- Go 1.26 以上

## 開発コマンド

| コマンド | 説明 |
|----------|------|
| `make build` | CLI バイナリを `bin/load` にビルドする |
| `make test` | ユニットテストを実行する |
| `make vet` | 静的解析を実行する |
| `make run` | デフォルト設定 (CPU 80%, Mem 70%, 30s) でローカル実行する |
| `make clean` | ビルド成果物 (`bin/`) を削除する |

```bash
# ビルド
make build

# テスト
make test

# ローカル実行（CPU 80%、メモリ 70%、30秒）
make run
```
```

## リリース

`v` から始まるタグを push すると GitHub Actions が自動でリリースを作成する。
リリースノートにはインストール方法・使い方・自動生成チェンジログが含まれる。

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 内部構造

| パッケージ | 役割 |
|-----------|------|
| `load` (ルート) | `Load` 構造体・Functional Options パターン・ゴルーチン管理 |
| `cpu` | gopsutil で CPU 時間を取得し、ビジーループで指定割合まで負荷を発生させる |
| `memory` | gopsutil で物理メモリ情報を取得し、差分サイズ分のバッファを確保する |
| `cmd/load` | CLI エントリーポイント |

## ライセンス

MIT

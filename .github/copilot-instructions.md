# go-load — Copilot Instructions

## プロジェクト概要

CPU・メモリに対して指定した割合の負荷を一定時間かけるライブラリおよびCLIツール。
負荷試験・スケールテスト・監視アラートの動作確認用途を想定している。

## パッケージ構成

```
go-load/
├── load.go            # Load 構造体・Functional Options・Run()
├── load_test.go       # テスト
├── cpu/
│   └── cpu.go         # CPU 負荷（duty-cycle 方式、外部依存なし）
├── memory/
│   └── memory.go      # メモリ負荷（/proc/meminfo 直読み、外部依存なし）
└── cmd/load/
    └── main.go        # CLI エントリーポイント
```

## コーディングルール

- Go 1.26 以上の構文・機能を使用する
- 外部依存ライブラリは持たない（標準ライブラリのみ）
- `cpu` パッケージ: `runtime` と `time` のみ使用
- `memory` パッケージ: `bufio`/`os`/`strconv`/`strings`/`runtime` のみ使用（`/proc/meminfo` を直接読む）
- エラーは呼び出し元に返す。内部パッケージ（`cpu`, `memory`）でのログ出力は禁止
- `Load.Run()` は `defer cancel()` で必ず context をキャンセルしてゴルーチンを終了させる
- Functional Options パターン（`OptionXxx`）でパラメータを追加する
- `Loader` インターフェース（`Load(context.Context) error`）に準拠した実装を追加すること

## 変更を加える際の注意点

### CPU (cpu/cpu.go)

- duty-cycle 方式を採用している。`windowSize = 100ms` のうち `per%` だけビジーループ、残りをスリープする
- `runtime.NumCPU()` 個のゴルーチンを起動し、各コアで並列に duty-cycle を回す
- gopsutil の `cpu.Times()` は使用しない（累積時間ベースなので精度が出ないため）

### メモリ (memory/memory.go)

- 割り当て前に `runtime.GC()` + `debug.FreeOSMemory()` を実行して正確なベースラインを取る
- 実使用量の基準は `vm.Total - vm.Available`（監視ツール `free`/`htop` の表示と一致する）
- `vm.Used` は buff/cache の計算方法が環境によって異なるため使用しない
- バッファへの書き込みはページ単位（4096 バイト間隔）で行い、物理メモリへのコミットと処理速度を両立する
- `runtime.KeepAlive(buf)` で GC によるバッファ解放を防ぐ
- `target <= currentUsed` の場合は即座に待機に入る（負の差分によるパニック防止）

### 新しい負荷種別を追加する場合

- `<種別>/` サブパッケージを作り `Loader` インターフェース（`Load(context.Context) error`）を実装する
- CLI フラグを追加した場合は README.md のフラグ一覧テーブルも更新する

## テスト・検証コマンド

```bash
# ビルド確認
make build
# または
go build ./...

# テスト実行
make test
# または
go test ./...

# 静的解析
make vet
# または
go vet ./...

# ローカル実行（CPU 80%、メモリ 70%、30秒）
make run

# ビルド成果物の削除
make clean
```

## Makefile ターゲット一覧

| ターゲット | 説明 |
|-----------|------|
| `make build` | CLI バイナリを `bin/load` にビルドする |
| `make test` | ユニットテストを実行する |
| `make vet` | 静的解析を実行する |
| `make run` | デフォルト設定 (CPU 80%, Mem 70%, 30s) でローカル実行する |
| `make clean` | ビルド成果物 (`bin/`) を削除する |

## 依存更新

外部依存なし。標準ライブラリのみ使用のため `go mod tidy` のみで十分。

```bash
go mod tidy
```

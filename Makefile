BIN := load

## help: このヘルプを表示する
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'

## build: CLI バイナリをビルドする
build:
	go build -o bin/$(BIN) ./cmd/load

## test: ユニットテストを実行する
test:
	go test ./...

## vet: 静的解析を実行する
vet:
	go vet ./...

## run: デフォルト設定 (CPU 80%, Mem 70%, 30s) でローカル実行する
run:
	go run ./cmd/load -duration 30s -cpu 80 -mem 70

## clean: ビルド成果物を削除する
clean:
	rm -rf bin/

.PHONY: help build test vet run clean

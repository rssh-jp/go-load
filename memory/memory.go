package memory

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

const pageSize = 4096

type Instance struct {
	per float64
}

func New(per float64) *Instance {
	return &Instance{per: per}
}

// virtualMemoryStat は /proc/meminfo から取得したメモリ情報を保持する
type virtualMemoryStat struct {
	total     uint64 // MemTotal (bytes)
	available uint64 // MemAvailable (bytes)
}

// readVirtualMemory は /proc/meminfo を読み取りメモリ情報を返す
func readVirtualMemory() (*virtualMemoryStat, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("open /proc/meminfo: %w", err)
	}
	defer f.Close()

	stat := &virtualMemoryStat{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		// 値は kB 単位なので 1024 を掛けてバイトに変換
		val, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		switch strings.TrimSuffix(fields[0], ":") {
		case "MemTotal":
			stat.total = val * 1024
		case "MemAvailable":
			stat.available = val * 1024
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read /proc/meminfo: %w", err)
	}
	if stat.total == 0 {
		return nil, fmt.Errorf("/proc/meminfo: MemTotal not found")
	}
	return stat, nil
}

func (inst Instance) Load(ctx context.Context) error {
	// GC を実行して正確な使用量を計測する
	runtime.GC()
	debug.FreeOSMemory()

	vm, err := readVirtualMemory()
	if err != nil {
		return err
	}

	// Total - Available を「実使用量」とする（監視ツールの表示と一致）
	currentUsed := vm.total - vm.available
	target := uint64(float64(vm.total) * inst.per / 100)

	if target <= currentUsed {
		// 既に目標使用量を超えているので追加割り当て不要
		<-ctx.Done()
		return nil
	}
	diff := target - currentUsed

	buf := make([]byte, diff)
	// ページ単位でアクセスして物理メモリにコミットさせる
	for i := 0; i < len(buf); i += pageSize {
		buf[i] = 1
	}
	if len(buf) > 0 {
		buf[len(buf)-1] = 1
	}

	// buf を GC されないよう保持したまま終了を待つ
	<-ctx.Done()
	runtime.KeepAlive(buf)
	return nil
}

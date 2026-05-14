package memory

import (
	"context"
	"runtime"
	"runtime/debug"

	"github.com/shirou/gopsutil/v3/mem"
)

const pageSize = 4096

type Instance struct {
	per float64
}

func New(per float64) *Instance {
	return &Instance{per: per}
}

func (inst Instance) Load(ctx context.Context) error {
	// GC を実行して正確な使用量を計測する
	runtime.GC()
	debug.FreeOSMemory()

	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	// Total - Available を「実使用量」とする（監視ツールの表示と一致）
	currentUsed := vm.Total - vm.Available
	target := uint64(float64(vm.Total) * inst.per / 100)

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

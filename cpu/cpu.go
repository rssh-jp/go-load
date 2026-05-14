package cpu

import (
	"context"
	"runtime"
	"time"
)

// windowSize は duty-cycle の1周期
const windowSize = 100 * time.Millisecond

type Instance struct {
	per float64
}

func New(per float64) *Instance {
	return &Instance{per: per}
}

func (inst Instance) Load(ctx context.Context) error {
	numCPU := runtime.NumCPU()
	done := make(chan struct{}, numCPU)

	for i := 0; i < numCPU; i++ {
		go func() {
			loadCore(ctx, inst.per)
			done <- struct{}{}
		}()
	}

	for i := 0; i < numCPU; i++ {
		<-done
	}
	return nil
}

// loadCore は 1コア分のデューティサイクルを実行する。
// per% の時間だけビジーループ、残り (100-per)% はスリープ。
func loadCore(ctx context.Context, per float64) {
	busy := time.Duration(float64(windowSize) * per / 100)
	idle := windowSize - busy

	for {
		// ビジーフェーズ
		if busy > 0 {
			end := time.Now().Add(busy)
			for time.Now().Before(end) {
				// 意図的なビジーループ
			}
		}

		// アイドルフェーズ
		if idle <= 0 {
			select {
			case <-ctx.Done():
				return
			default:
			}
			continue
		}

		timer := time.NewTimer(idle)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

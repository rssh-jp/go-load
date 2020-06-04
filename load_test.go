package load

import (
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	l := New(OptionMemoryPercentage(100), OptionCPUPercentage(50), OptionDuration(time.Second*5))
	l.Run()
}

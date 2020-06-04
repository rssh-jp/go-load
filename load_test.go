package load

import (
	"testing"
	"time"
)

func TestSuccess(t *testing.T) {
	l := New(OptionMemoryPercentage(50), OptionCPUPercentage(50), OptionDuration(time.Second*5))
	l.Run()
}

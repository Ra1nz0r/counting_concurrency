package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ra1nz0r/counting_concurrency/internal/dye"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	chanInT := make(chan int64)
	var num int64 = 1
	var compareNum int64

	go func() {
		defer close(chanInT)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				chanInT <- num
				compareNum = num
				num++
			}
		}
	}()

	chanOutT := make(chan int64)

	go Worker(chanInT, chanOutT)

	require.NotEmpty(t, <-chanOutT,
		fmt.Sprintf("%sChannel empty.%s", dye.Red, dye.Reset))

	var count int64
	for v := range chanOutT {
		count = v
	}

	assert.Equal(t, compareNum, count,
		fmt.Sprintf("%sWrong count. Expect: %d, got: %d.%s", dye.Red, compareNum, count, dye.Reset))
}

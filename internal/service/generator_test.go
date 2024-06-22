package service

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ra1nz0r/counting_concurrency/internal/dye"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var testCntTOT int64
	var testSumTOT int64

	chInT := GenNum(ctx, &testSumTOT, &testCntTOT)

	require.NotEmpty(t, <-chInT,
		fmt.Sprintf("%sChannel empty.%s", dye.Red, dye.Reset))

	var expCount int64 = 1
	var expSum int64 = 1

	for {
		value, ok := <-chInT
		if !ok {
			break
		}
		atomic.AddInt64(&expSum, value)
		atomic.AddInt64(&expCount, 1)
		time.Sleep(25 * time.Microsecond)
	}

	require.NotEqual(t, expCount, expSum,
		fmt.Sprintf("%s`Generator` does not increment the value.%s", dye.Red, dye.Reset))
	assert.Equal(t, expCount, testCntTOT,
		fmt.Sprintf("%sWrong count. Expect: %d, got: %d.%s", dye.Red, expCount, testCntTOT, dye.Reset))
	assert.Equal(t, expSum, testSumTOT,
		fmt.Sprintf("%sWrong sum. Expect: %d, got: %d.%s", dye.Red, expSum, testSumTOT, dye.Reset))
}

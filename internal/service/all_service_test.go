package service

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/ra1nz0r/counting_concurrency/internal/dye"

	"github.com/stretchr/testify/assert"
)

func TestAllService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var inSumT int64
	var inCntT int64

	chIn := GenNum(ctx, &inSumT, &inCntT)

	NumOut := runtime.NumCPU()
	outs := GenChanSliceWithNum(NumOut, chIn)

	chOut, amounts := SendNumInResChan(NumOut, outs)

	sumTestTOT, cntTestTOT := ReadFromResChan(chOut)

	assert.Equal(t, inCntT, cntTestTOT,
		fmt.Sprintf("%sCount numbers is not equal.%s", dye.Red, dye.Reset))
	assert.Equal(t, inSumT, sumTestTOT,
		fmt.Sprintf("%sSum numbers is not equal.%s", dye.Red, dye.Reset))

	for _, v := range amounts {
		inCntT -= v
	}
	assert.Zero(t, inCntT,
		fmt.Sprintf("%sThe division of numbers by channel is incorrect.%s", dye.Red, dye.Reset))
}

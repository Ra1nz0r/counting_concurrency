package main

import (
	"context"
	"runtime"
	"time"

	"github.com/ra1nz0r/counting_concurrency/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var inSum int64 // сумма сгенерированных чисел
	var inCnt int64 // количество сгенерированных чисел

	chIn := service.GenNum(ctx, &inSum, &inCnt)

	NumOut := runtime.NumCPU() // количество обрабатывающих горутин и каналов
	outs := service.GenChanSliceWithNum(NumOut, chIn)

	chOut, amounts := service.SendNumInResChan(NumOut, outs)

	sumTOT, cntTOT := service.ReadFromResChan(chOut)

	service.CheckResult(amounts, inSum, inCnt, sumTOT, cntTOT)
}

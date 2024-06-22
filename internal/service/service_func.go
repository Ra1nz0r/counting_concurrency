package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Generator генерирует последовательность чисел 1,2,3 и т.д. и
// отправляет их в канал ch. При этом после записи в канал для каждого числа
// вызывается функция fn. Она служит для подсчёта количества и суммы
// сгенерированных чисел.
func Generator(ctx context.Context, ch chan<- int64, fn func(int64)) {
	var num int64 = 1
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ch <- num
			fn(num)
			num++
		}
	}
}

// Worker читает число из канала in и пишет его в канал out.
func Worker(in <-chan int64, out chan<- int64) {
	defer close(out)
	for {
		value, ok := <-in
		if !ok {
			break
		}
		out <- value
		time.Sleep(time.Millisecond)
	}
}

// Генерирует числа и считает их количество и сумму.
func GenNum(ctx context.Context, numSum, numCnt *int64) chan int64 {
	chInput := make(chan int64)
	go Generator(ctx, chInput, func(i int64) {
		atomic.AddInt64((*int64)(numSum), i)
		atomic.AddInt64((*int64)(numCnt), 1)
	})
	return chInput
}

// Возвращает слайс каналов, с количеством записанных чисел в каждый канал из входного канала.
func GenChanSliceWithNum(numOutput int, chInput chan int64) []chan int64 {
	// outsChanSlice — слайс каналов, куда будут записываться числа из chIn
	outsChanSlice := make([]chan int64, numOutput)
	for i := 0; i < numOutput; i++ {
		// создаём каналы и для каждого из них вызываем горутину Worker
		outsChanSlice[i] = make(chan int64)
		go Worker(chInput, outsChanSlice[i])
	}
	return outsChanSlice
}

// Отправляет числа из слайса каналов в результирующий канал, со статистикой по отработанным горутинам.
func SendNumInResChan(numOutput int, outsChanSlice []chan int64) (chan int64, []int64) {
	// gorutStat — слайс, в который собирается статистика по горутинам
	gorutStat := make([]int64, numOutput)
	// chanRes — канал, в который будут отправляться числа из горутин `gorutStat[i]`
	chanRes := make(chan int64, numOutput)

	var wg sync.WaitGroup

	// Передаем числа из каналов gorutStat в результирующий канал
	for i := 0; i < len(outsChanSlice); i++ {
		wg.Add(1)
		go func(nextChan <-chan int64, numChan int64) {
			defer wg.Done()
			for v := range nextChan {
				chanRes <- v
				gorutStat[numChan]++
			}
		}(outsChanSlice[i], int64(i))
	}

	go func() {
		wg.Wait()
		close(chanRes)
	}()
	return chanRes, gorutStat
}

// Читаем числа из результирующего канала
func ReadFromResChan(chanRes chan int64) (checkSum, checkCnt int64) {
	for v := range chanRes {
		checkCnt++
		checkSum += v
	}
	return
}

// Проверка результата и вывод данных по числам, каналам.
func CheckResult(gorutStat []int64, numSum, numCnt, checkSum, checkCnt int64) {
	fmt.Println("Количество чисел", numCnt, checkCnt)
	fmt.Println("Сумма чисел", numSum, checkSum)
	fmt.Println("Разбивка по каналам", gorutStat)

	if numSum != checkSum {
		log.Fatalf("Ошибка: суммы чисел не равны: %d != %d\n", numSum, checkSum)
	}
	if numCnt != checkCnt {
		log.Fatalf("Ошибка: количество чисел не равно: %d != %d\n", numCnt, checkCnt)
	}
	for _, v := range gorutStat {
		numCnt -= v
	}
	if numCnt != 0 {
		log.Fatalf("Ошибка: разделение чисел по каналам неверное\n")
	}
}

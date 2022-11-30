// Напишите код, в котором имеются два канала сообщений из целых чисел так,
// чтобы приём сообщений всегда приводил к блокировке.
// Приёмом сообщений из обоих каналов будет заниматься главная горутина.
//
// Сделайте так, чтобы во время такого «бесконечного ожидания» сообщений
// выполнялась фоновая работа в виде вывода текущего времени в консоль
package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	// Инициализация каналов ch1 и ch2
	ch1, ch2 := make(chan int), make(chan int)

	//ctx и cancel Для уведомления горутины о необходимости завершить свою работу.
	ctx, cancel := context.WithCancel(context.TODO())
	wg := sync.WaitGroup{}

	wg.Add(2)
	go startTicker(ctx, ticker)
	go startWriter(ctx, &wg, ch1, time.Second*4, 5)
	go startWriter(ctx, &wg, ch2, time.Second*6, 5)

	go startWaiter(ctx, ch1, ch2)
	wg.Wait()
	cancel()
}

// функция startWriter принимает значения в правильном порядке
// и выводит на экран (в консоль).
func startWriter(ctx context.Context, wg *sync.WaitGroup, ch chan int, d time.Duration, routines int) {
	for i := 0; i < routines; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(d)
			ch <- i
		}
	}
	wg.Done()
}

// Фуекция startWaiter по порядку отсылает числа от в каналы ch1 и ch2 .
func startWaiter(ctx context.Context, ch1, ch2 chan int) {
	for {
		select {
		case v := <-ch1:
			fmt.Printf("получено от ch1: %d\n", v)
		case v := <-ch2:
			fmt.Printf("получено от ch2: %d\n", v)
		case <-ctx.Done():
			fmt.Println("done")
			os.Exit(0)
		}
	}
}

// Фуекция startTicker выводит текущие время в консоль.
func startTicker(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t := <-ticker.C
			outputMessage := []byte("Время: ")
			outputMessage = t.AppendFormat(outputMessage, "15:04:05")
			fmt.Println(string(outputMessage))
		}
	}
}

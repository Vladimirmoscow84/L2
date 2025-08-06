/*Реализовать функцию, которая будет объединять один или более каналов done (каналов сигнала завершения) в один. Возвращаемый канал должен закрываться, как только закроется любой из исходных каналов.
Сигнатура функции может быть такой:

var or func(channels ...<-chan interface{}) &lt;-chan interface{}
Пример использования функции:

sig := func(after time.Duration) <-chan interface{} {
   c := make(chan interface{})
   go func() {
      defer close(c)
      time.Sleep(after)
   }()
   return c
}

start := time.Now()
<-or(
   sig(2*time.Hour),
   sig(5*time.Minute),
   sig(1*time.Second),
   sig(1*time.Hour),
   sig(1*time.Minute),
)
fmt.Printf("done after %v", time.Since(start))
В этом примере канал, возвращённый or(...), закроется через ~1 секунду, потому что самый короткий канал sig(1*time.Second) закроется первым. Ваша реализация or должна уметь принимать на вход произвольное число каналов и завершаться при сигнале на любом из них.

Подсказка: используйте select в бесконечном цикле для чтения из всех каналов одновременно, либо рекурсивно объединяйте каналы попарно.*/

package main

import (
	"fmt"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	//время начала работы
	start := time.Now()

	//прием в функцию maybe  нескольких функций, с открытыми каналами (с определенным таймаутом)
	<-maybe(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Minute),
		sig(5*time.Second),
		sig(999*time.Millisecond),
		sig(5555*time.Millisecond),
		sig(2*time.Second),
	)
	//при закрытиии одного из каналов закрывается общий канал, разблокировав главную горутину
	fmt.Printf("done after %v", time.Since(start))
}

func maybe(chans ...<-chan interface{}) <-chan interface{} {

	chDone := make(chan interface{})
	chOut := make(chan interface{})
	for i := range chans {

		go func(ch <-chan interface{}) {

			select {
			case signal := <-ch:
				close(chDone)
				chOut <- signal
			case <-chDone:
				return
			}

		}(chans[i])
	}
	<-chDone
	return chOut
}

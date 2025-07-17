/*
Что выведет программа?

Объяснить работу конвейера с использованием select.
Ответ: Программа выведет данные, прочитанные из канала с в произвольном порядке, ранее записанные в него из каналов a и b .
в функции merge объединение двух канало в один. В горутине в бесконечном цикле есть оператр select, который имеет два кейса,
в которых порисходит чтение из каналов а и б . Происходит проверка, поступают ли данные из канала, если данне есть, то  ok(true),пишем эти данне в канал С, если данные не поступают (канал закрыт), то !ok(false), в этом случаe присваеваем каналу значеие nil -  канал становится янедействителтным(мы с ним больше не рабтаем).
После закрытия двух каналов , из котрых читаем (им присвоен nil) по условию мы закрываем канал, в который пишем, А в главной горутине завершается чтение из этого канала и происходит выход из программы
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

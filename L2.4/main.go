/*то выведет программа?

Объяснить вывод программы.

func main() {
  ch := make(chan int)
  go func() {
    for i := 0; i &lt; 10; i++ {
    ch &lt;- i
  }
}()
  for n := range ch {
    println(n)
  }
}*/

/*выведутся цифры от 0 до 9 и случится deadlock, так как канал ch не закрылся. в главная горутина, после прочтения всех данных будет все ще ждать данные из канала
Для корректного выполенния программы необходимо закрыть канал в горутине после цикла*/

package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		//close(ch)
	}()
	for n := range ch {
		println(n)
	}
}

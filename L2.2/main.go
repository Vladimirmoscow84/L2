/*Что выведет программа?

Объяснить порядок выполнения defer функций и итоговый вывод.

package main

import "fmt"

func test() (x int) {
  defer func() {
    x++
  }()
  x = 1
  return
}

func anotherTest() int {
  var x int
  defer func() {
    x++
  }()
  x = 1
  return x
}

func main() {
  fmt.Println(test())
  fmt.Println(anotherTest())
}*/

package main

import "fmt"

/*ответ 2 и 1
в функции test() (x int) перепменная х именована в сигнатуре функции и передана в анонимную функцию, после чего по коду программы ей присваивается значение 1, а после оператора return (завершение функции) срабатывает defer в анонимной функции происходит ее увеличение на 1, а затем происходит возврат этой переменной со сзначением 2(1+1)
в функции anotherTest() int переменная х объявлена внутри функции, инициализированна х=1 и возвращена с этим значением, а оператор defer и вычисления в анонимной функци происходят после возвращения значения х==1*/

func test() (x int) {
	defer func() {
		x++

	}()
	x = 1
	return x
}

func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())

}

//Объяснить вывод программы.

/*Программа выведет error,
err является интерфейсом типа error и он является nil, так как у него value nil и type nil, в функции тест ему присвоили type customError value nil, так как теперь переменная интерфесного типа err не является nil, то утвержедение err!!=nil true и мы проваливаемся во внутрь оператора  */

package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

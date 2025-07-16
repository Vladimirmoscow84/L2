/*
Что выведет программа?
Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

package main

import (

	"fmt"
	"os"

)

	func Foo() error {
	  var err *os.PathError = nil
	  return err
	}

	func main() {
	  err := Foo()
	  fmt.Println(err)
	  fmt.Println(err == nil)
	}
*/

/*
Ответ:
программа вывведет nil и false,
интерфейс является нил, когда у него тип и значение равны нил. в мэйн err является экземпляром интерфейса с типом os.PathError и значением нил
Первый принтлн выводит значение(value), а второй неверное удверждение, так как данная переменная интерфейсного типа не нил
Пустой интерефейс не содержит методов, поэтому ему соответствует любой тип
*/
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}

/*Распаковка строки
Написать функцию Go, осуществляющую примитивную распаковку строки, содержащей повторяющиеся символы/руны.

Примеры работы функции:

Вход: "a4bc2d5e"
Выход: "aaaabccddddde"

Вход: "abcd"
Выход: "abcd" (нет цифр — ничего не меняется)

Вход: "45"
Выход: "" (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)

Вход: ""
Выход: "" (пустая строка -> пустая строка)

Дополнительное задание
Поддерживать escape-последовательности вида \:

Вход: "qwe\4\5"
Выход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)

Вход: "qwe\45"
Выход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)

Требования к реализации
Функция должна корректно обрабатывать ошибочные случаи (возвращать ошибку, например, через error), и проходить unit-тесты.

Код должен быть статически анализируем (vet, golint).*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	a := "\\4a3as4ffg\\42k3\\54"
	fmt.Println("")
	fmt.Println(getUnpacking(a))

}

func getUnpacking(data string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var previos, result string
	var count int
	previosCharIsDigit := true
	previosCharIsSlash := true
	previosCharIsLetter := true

	for _, v := range data {
		if unicode.IsDigit(v) && previosCharIsDigit {
			err := errors.New("некорректная строка")
			return "", err
		}
		if v == '\\' {
			result += previos
			previos = ""
			previosCharIsSlash = true
			previosCharIsLetter = true
			previosCharIsDigit = false
			continue
		}
		if unicode.IsLetter(v) {
			result += previos
			previos = string(v)
			previosCharIsDigit = false
			previosCharIsLetter = true
			previosCharIsSlash = false
			continue
		}

		if unicode.IsDigit(v) && previosCharIsLetter {
			if !previosCharIsSlash {
				count, _ = strconv.Atoi(string(v))
				result += strings.Repeat(previos, count)
				previos = ""
				previosCharIsLetter = false
				previosCharIsDigit = true
				previosCharIsSlash = false
				continue
			}
		}

		if unicode.IsDigit(v) && previosCharIsSlash {
			result += previos
			previos = string(v)
			previosCharIsDigit = false
			previosCharIsLetter = true
			previosCharIsSlash = false
		}

	}

	result += previos
	fmt.Println(result)
	return result, nil
}

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
	a := "\\4"
	fmt.Println(getUnpacking(a))

}

func getUnpacking(data string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var previos, result string
	var count int
	previosCharIsLetter := true
	previosCharIsSlash := true

	for i, v := range data {
		if (unicode.IsDigit(v) && i == 0) || (unicode.IsDigit(v) && !previosCharIsLetter && !previosCharIsSlash) {
			err := errors.New("некорректная строка")
			return "", err
		}
		if unicode.IsLetter(v) {
			result += previos
			previos = string(v)
			previosCharIsLetter = true
			continue
		}

		if unicode.IsDigit(v) && previosCharIsLetter {
			count, _ = strconv.Atoi(string(v))
			result += strings.Repeat(previos, count)
			previos = ""
			previosCharIsLetter = false
			continue
		}
		if v == '\\' {
			result += previos
			previos = ""
			previosCharIsSlash = true
			continue
		}
		if unicode.IsDigit(v) && previosCharIsSlash {
			previos = string(v)
			previosCharIsLetter = true
			previosCharIsSlash = false
		}

	}

	result += previos
	return result, nil
}

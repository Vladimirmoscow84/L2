/*Утилита cut
Реализовать утилиту, которая считывает входные данные (STDIN) и разбивает каждую строку по заданному разделителю, после чего выводит определённые поля (колонки).

Аналог команды cut с поддержкой флагов:

-f "fields" — указание номеров полей (колонок), которые нужно вывести. Номера через запятую, можно диапазоны.
Например: «-f 1,3-5» — вывести 1-й и с 3-го по 5-й столбцы.

-d "delimiter" — использовать другой разделитель (символ). По умолчанию разделитель — табуляция ('\t').

-s – (separated) только строки, содержащие разделитель. Если флаг указан, то строки без разделителя игнорируются (не выводятся).

Программа должна корректно парсить аргументы, поддерживать различные комбинации (например, несколько отдельных полей и диапазонов), учитывать, что номера полей могут выходить за границы (в таком случае эти поля просто игнорируются).

Стоит обратить внимание на эффективность при обработке больших файлов. Все стандартные требования по качеству кода и тестам также применимы.*/

package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

//структура для использования аргументов командной строки

type arguments struct {
	f     string
	d     string
	input string
	s     bool
}

func main() {

}

// parseFlags парсит арументы командной строки
func parseFlags() (arguments, error) {

	var f, d string
	var s bool

	//определение флагов
	flag.StringVar(&f, "f", "", "\"fields\" -указание номеров полей (колонок)")
	flag.StringVar(&d, "d", "\t", "\"delimiter\" -использовать другой разделитель (символ)")
	flag.BoolVar(&s, "s", false, "\"separated\"- строки с разделителем")

	//парсинг флагов  в переменнные
	flag.Parse()
	input := flag.Arg(0)

	//проверка введенной строки на валидность
	if input == "" {
		return arguments{}, fmt.Errorf("строка невалидна")
	}

	//создается экземпляр структуры с аргументами командной строки
	cmdArgs := arguments{
		f:     f,
		d:     d,
		s:     s,
		input: input,
	}
	return cmdArgs, nil

}

// cmdDefinSelect - определяет и выбирает колонки из аргумента командной строки
func cmdDefinSelect(cmdArgs arguments, rows []string) (answer []int, err error) {

	//конвертация строкового занчения числа в int
	var num int

	//преобразвоание строки в слайс для конвертации в int
	nums := strings.Split(cmdArgs.f, ",")

	//преобразование строковых чисел в int
	for _, v := range nums {
		if strings.Contains(v, "-") && len(v) > 1 && strings.Count(v, "-") == 1 {
			if strings.Index(v, "-") == 0 {
				num, err = strconv.Atoi(strings.TrimPrefix(v, "-"))
				if err != nil {
					return nil, fmt.Errorf("неверный номер колонки")
				}
				for i := 0; i< num; i++ {
					if i>len(rows)-1{
						return answer, nil
					}
					answer = append(answer, i)
				}
				continue

			}else if  {
				
			}
		}
	}

}

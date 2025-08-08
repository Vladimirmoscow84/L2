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
	args, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	//в случае заданного ключа только с разделителем, а в строке нет разделителя - выход из программы
	if args.s {
		if !strings.Contains(args.input, args.d) {
			return
		}
	}
	//разделение строкина части по разделителю
	sepString := strings.Split(args.input, args.d)
	//проверка ключа на выбор полей и что строка разделена
	if args.f != "" && len(sepString) > 1 {
		nums, err := cmdDefinSelect(args, sepString)
		if err != nil {
			fmt.Println(err)
			return
		}

		//вывод колонок в соответстви с условием
		for i, v := range nums {
			if i == len(nums)-1 {
				fmt.Println(sepString[v])
				return
			}
			fmt.Print(sepString[v] + args.d)
		}

	}
	//в случае несрабатывания условий проверки выводится исходная строка
	fmt.Println(args.input)

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
				for i := 0; i < num; i++ {
					if i > len(rows)-1 {
						return answer, nil
					}
					answer = append(answer, i)
				}
				continue
				// Если "-" является последнием символом после числа,
				// то берём все числа, пока последнее не станет равно номеру последней колонки,
				// конвертируем их в числовое представление и добавляем в массив

			} else if strings.Index(v, "-") == 0 {
				num, err = strconv.Atoi(strings.TrimSuffix(v, "-"))
				if err != nil {
					return nil, fmt.Errorf("неверный номер колонки")
				}
				for j := num - 1; j < len(rows); j++ {
					answer = append(answer, j)
				}
				continue
				// Если "-" находится между двумя числами,
				// берём все числа в этом промежутке и конвертируем их числовое представление,а затем добавляем в массив
			} else {
				twoNums := strings.Split(v, "-")
				one, err := strconv.Atoi(twoNums[0])
				if err != nil {
					return nil, err
				}
				two, err := strconv.Atoi(twoNums[1])
				if err != nil {
					return nil, err
				}
				for j := one - 1; j < two; j++ {
					if j > len(rows)-1 {
						return answer, nil
					}
					answer = append(answer, j)
				}
				continue

			}
		}
		//Конвертируем строковое представление в числовое, в случае ошибки возвращаем её
		num, err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("колонка с номерром %d не существует", num)
		}
		//полученное число добавляем в массив
		answer = append(answer, num-1)
	}
	return answer, nil

}

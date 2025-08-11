/*Утилита sort
Реализовать упрощённый аналог UNIX-утилиты sort (сортировка строк).

Программа должна читать строки (из файла или STDIN) и выводить их отсортированными.

Обязательные флаги (как в GNU sort):

-k N — сортировать по столбцу (колонке) №N (разделитель — табуляция по умолчанию).
Например, «sort -k 2» отсортирует строки по второму столбцу каждой строки.

-n — сортировать по числовому значению (строки интерпретируются как числа).

-r — сортировать в обратном порядке (reverse).

-u — не выводить повторяющиеся строки (только уникальные).

Дополнительные флаги:

-M — сортировать по названию месяца (Jan, Feb, ... Dec), т.е. распознавать специфический формат дат.

-b — игнорировать хвостовые пробелы (trailing blanks).

-c — проверить, отсортированы ли данные; если нет, вывести сообщение об этом.

-h — сортировать по числовому значению с учётом суффиксов (например, К = килобайт, М = мегабайт — человекочитаемые размеры).

Программа должна корректно обрабатывать комбинации флагов (например, -nr — числовая сортировка в обратном порядке, и т.д.).

Необходимо предусмотреть эффективную обработку больших файлов.

Код должен проходить все тесты, а также проверки go vet и golint (понимание, что требуются надлежащие комментарии, имена и структура программы).*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type arguments struct {
	k                   int
	n, r, u, m, b, c, h bool
	input, output       string
}

// parseFlags - распарсивает аргументы ккомандной строки
func parseFlags() (arguments, error) {

	//переменнные для хранения значений флагов
	var k int
	var n, r, u, m, b, c, h bool

	//определем флаги и парсим их
	flag.BoolVar(&n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&r, "r", false, "сортировать в обратном порядке (reverse)")
	flag.BoolVar(&u, "u", false, "не выводить повторяющиеся строки (только уникальные)")
	flag.BoolVar(&m, "m", false, "сортировать по названию месяца")
	flag.BoolVar(&b, "b", false, "игнорировать хвостовые пробелы (trailing blanks)")
	flag.BoolVar(&c, "c", false, "проверить, отсортированы ли данные; если нет, вывести сообщение об этом")
	flag.BoolVar(&h, "h", false, "сортировать по числовому значению с учётом суффиксов")
	flag.IntVar(&k, "k", 0, "колонка сортировки")
	flag.Parse()

	input := flag.Arg(0)
	output := flag.Arg(1)

	//проверка на валидность имен файлов
	if input == "" && output == "" {
		return arguments{}, fmt.Errorf("отсутствие имен файлов")
	}

	//инициализация экземпляра структуры с элеентами командной строки
	args := arguments{
		k:      k,
		n:      n,
		r:      r,
		u:      u,
		m:      m,
		b:      b,
		c:      c,
		h:      h,
		input:  input,
		output: output,
	}
	return args, nil

}

// readFile - читает данные перед сортировкой
func readFile(args arguments) ([][]string, error) {
	//инициализируем матрицу для хранения данных
	data := make([][]string, 0)
	//открываем файл
	file, err := os.Open(args.input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//сканируем данные из файла
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		l := scanner.Text()
		//если передан агругмент на игнорирование хвостовых пробелов, то обрезаем
		if args.b {
			l = strings.TrimSuffix(l, " ")
		}
		//разрезаем строку на слайс строк и добавляем в матрицу данных
		elements := strings.Split(l, " ")
		data = append(data, elements)
	}
	return data, nil

}

// writeFile - записывает отсортированные данные в файл
func writeFile(data [][]string, args arguments) error {
	//создание файла для записи
	file, err := os.Create(args.output)
	if err != nil {
		return err
	}
	defer file.Close()

	//проверка повторения строк (создается мапа для проверки по ключу)
	repStrings := make(map[string]struct{}, len(data))
	//создаем слайс строк, для добавления туда отсортированных строк
	srtStrings := make([]string, len(data))
	for _, value := range data {
		l := strings.Join(value, " ")
		//проверка на наличие в мапе (на повтор)
		if args.u {
			if _, ok := repStrings[l]; ok {
				continue
			}
		}
		//если нет такой строки, то добавляем в мапу
		repStrings[l] = struct{}{}

		//добавляем в слайс новую отсортированную строку
		srtStrings = append(srtStrings, l)
	}
	//отсортированные данные записываем в файл
	_, err = file.WriteString(strings.Join(srtStrings, "\n"))
	if err != nil {
		return err
	}
	return nil

}

// selectSortFunc - выбор алгоритма сортировки данных
func selectSortFunc(data [][]string, args arguments) func(i, j int) bool {
	// Объявляем функцию для сортировки данных
	var sortFunc func(i, j int) bool
	// В зависимости от переданных аргументов присваиваем ей алгоритм сравнивания элементов
	switch {
	case args.n:
		// Сравнение данные по числовым значениям
		sortFunc = func(i, j int) bool {
			firstElem, _ := strconv.ParseFloat(getElement(data, i, args.k), 64)
			secElem, _ := strconv.ParseFloat(getElement(data, j, args.k), 64)
			if args.r {
				return firstElem > secElem
			}
			return firstElem < secElem
		}
		// Сравнение данных по названию месяца
	case args.m:
		sortFunc = func(i, j int) bool {
			firstElem, _ := getMonth(getElement(data, i, args.k))
			secElem, _ := getMonth(getElement(data, j, args.k))
			if args.r {
				return firstElem > secElem
			}
			return firstElem < secElem
		}
		// сравнение данных по суффиксам строк, если номер переданной колонки для сортировки
		// превышает длину строки, то размер элемента равен нулю
	case args.h:
		var firstElem, secElem int
		sortFunc = func(i, j int) bool {
			if args.k >= len(data[i]) {
				firstElem = 0
			} else {
				firstElem = getLen(data[i][args.k:])
			}
			if args.k >= len(data[j]) {
				secElem = 0
			} else {
				secElem = getLen(data[j][args.k:])
			}
			if args.r {
				return firstElem > secElem
			}
			return firstElem < secElem
		}
		// сравнение данных по размеру
	default:
		sortFunc = func(i, j int) bool {
			firstElem := getElement(data, i, args.k)
			secElem := getElement(data, j, args.k)
			if args.r {
				return firstElem > secElem
			}
			return firstElem < secElem
		}
	}
	return sortFunc
}

// Функция получает элемент исходной строки в зависимости от заданной колонки
func getElement(data [][]string, i, k int) string {
	if k < len(data[i]) {
		return data[i][k]
	}
	return ""
}

// Функция получает номер месяца
func getMonth(month string) (time.Month, error) {
	if m, err := time.Parse("January", month); err == nil {
		return m.Month(), nil
	}
	if m, err := time.Parse("Jan", month); err == nil {
		return m.Month(), nil
	}
	if m, err := time.Parse("01", month); err == nil {
		return m.Month(), nil
	}
	if m, err := time.Parse("1", month); err == nil {
		return m.Month(), nil
	}
	return 0, fmt.Errorf("не удалось получить месяц")
}

// Функция получает длину суффикса строки
func getLen(str []string) int {
	var sumLen int
	for _, v := range str {
		sumLen += len(v)
	}
	return sumLen
}
func main() {
	// Задаём и парсим аргументы командной строки
	args, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Открываем и считываем данные из файла, названия которого указывается в командной строке
	data, err := readFile(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Проверка и приведение номера сортируемой колонки к нужному индексу
	if args.k < 1 {
		args.k = 0
	} else {
		args.k--
	}

	// Выбор функции сортировки в зависимости от выбранных аргументов
	sortFunc := selectSortFunc(data, args)
	// Выполнение проверки на сортированность исходных данных при выбранном соответствующем аргументе
	if args.c {
		if sort.SliceIsSorted(data, sortFunc) {
			fmt.Println("sorted")
			return
		}
		fmt.Println("unsorted")
		return
	}
	// Сортировка исходных данных c помощью функции, выбранной ранее
	sort.Slice(data, sortFunc)

	// Запись отсортированных данных в файл, указанный в аргументах командной строки
	if err = writeFile(data, args); err != nil {
		fmt.Println(err)
		return
	}

}

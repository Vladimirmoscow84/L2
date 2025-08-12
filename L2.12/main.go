/*Утилита grep
Реализовать утилиту фильтрации текстового потока (аналог команды grep).

Программа должна читать входной поток (STDIN или файл) и выводить строки, соответствующие заданному шаблону (подстроке или регулярному выражению).

Необходимо поддерживать следующие флаги:

-A N — после каждой найденной строки дополнительно вывести N строк после неё (контекст).

-B N — вывести N строк до каждой найденной строки.

-C N — вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N).

-c — выводить только то количество строк, что совпадающих с шаблоном (т.е. вместо самих строк — число).

-i — игнорировать регистр.

-v — инвертировать фильтр: выводить строки, не содержащие шаблон.

-F — воспринимать шаблон как фиксированную строку, а не регулярное выражение (т.е. выполнять точное совпадение подстроки).

-n — выводить номер строки перед каждой найденной строкой.

Программа должна поддерживать сочетания флагов (например, -C 2 -n -i – 2 строки контекста, вывод номеров, без учета регистра и т.д.).

Результат работы должен максимально соответствовать поведению команды UNIX grep.

Обязательно учесть пограничные случаи (начало/конец файла для контекста, повторяющиеся совпадения и пр.).

Код должен быть чистым, отформатированным (gofmt), работать без ситуаций гонки и успешно проходить golint.*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// структура для предстаавлени аргументов командной мтроки
type arguments struct {
	A, B, C, idx           int
	c, i, v, F, n          bool
	oldMatch, match, input string
}

// ParseFlags распарсивает аргументы командной строки
func parseFlags() (arguments, error) {
	//определение перпеменных для хранения значений флагов
	var count, i, v, f, n bool
	var a, b, c int

	// опредееление флагов командной строки и распарсивание их в перменнные
	flag.BoolVar(&count, "c", false, "\"count\" (количество строк)")
	flag.BoolVar(&i, "i", false, "\"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&v, "v", false, "\"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&f, "F", false, "\"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&n, "n", false, "\"line num\", напечатать номер строки")
	flag.IntVar(&a, "A", 0, "\"after\" печатать +N строк после совпадения")
	flag.IntVar(&b, "B", 0, "\"before\" печатать +N строк до совпадения")
	flag.IntVar(&c, "C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.Parse()
	match := flag.Arg(0)
	input := flag.Arg(1)

	//проверка на валидность имен для файлов чтения данных
	if input == "" && match == "" {
		return arguments{}, fmt.Errorf("имя файла или искомая строка отсутствуют")
	}

	//инициализация экземпляра структуры с аргументами командной строки и инндексом вхождения
	args := arguments{
		A:        a,
		B:        b,
		C:        c,
		c:        count,
		i:        i,
		v:        v,
		F:        f,
		n:        n,
		oldMatch: match,
		match:    match,
		input:    input,
		idx:      -1,
	}
	return args, nil
}

// openFile - открывает и читает файл, возвращая считанный массив, в соответствии с заданным флагом
func openFile(args *arguments) ([]string, error) {
	file, err := os.Open(args.input)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]string, 0)
	scanner := bufio.NewScanner(file)

	//при считывании данных определяем флаг i, если он есть, топриводим все данные к нижнему регситру
	if args.i {
		args.match = strings.ToLower(args.match)
		for scanner.Scan() {
			l := strings.ToLower(scanner.Text())
			data = append(data, l)
		}
		return data, nil
	}
	for scanner.Scan() {
		l := scanner.Text()
		data = append(data, l)
	}
	return data, nil

}

// fullCoicidenceString - ищет строку по полному совпадению
func fullCoicidenceString(args *arguments, data []string) {
	for i, val := range data {
		if val == args.match {
			args.idx = i
			break
		}
	}
}

// findPostString - ищет позицию строки в переданных данных
func findPostString(args *arguments, data []string) {
	for i, val := range data {
		if strings.Contains(val, args.match) {
			args.idx = i
			break
		}
	}
}

// countRepString- считает кол-во повторений строки в данных
func countRepString(args *arguments, data []string) int {
	count := 0
	for _, val := range data {
		if strings.Contains(val, args.match) {
			count++
		}
	}
	return count
}

// delString - удаляет строку из данных
func delString(args *arguments, data []string) {
	//если есть флаг -F
	if args.F {
		for i, val := range data {
			if strings.Contains(val, args.match) {
				data = append(data[:i], data[i+1:]...)
			}
		}
		//если нет флага -F
	} else {
		for i, val := range data {
			if strings.Contains(val, args.match) {
				data[i] = strings.ReplaceAll(val, args.match, "")
			}
		}
	}
	fmt.Println()
	fmt.Printf("-v:\n\tДанные после удаления строки '%s':\n%s\n", args.oldMatch, strings.Join(data, "\n"))
	fmt.Println()

}

func main() {

}

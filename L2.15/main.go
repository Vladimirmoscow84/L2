/*Minishell: взаимодействие с ОС
Необходимо реализовать собственный простейший Unix shell.

Требования
Ваш интерпретатор командной строки должен поддерживать:

Встроенные команды:
– cd <path> – смена текущей директории.
– pwd – вывод текущей директории.
– echo <args> – вывод аргументов.
– kill <pid> – послать сигнал завершения процессу с заданным PID.
– ps – вывести список запущенных процессов.

Запуск внешних команд через exec (с помощью системных вызовов fork/exec либо стандартных функций os/exec).

Конвейеры (pipelines): возможность объединять команды через |, чтобы вывод одной команды направлять на ввод следующей (как в обычном shell).

Например: ps | grep myprocess | wc -l.

Обработку завершения: при нажатии Ctrl+D (EOF) шелл должен завершаться; Ctrl+C — прерывание текущей запущенной команды, но без закрыватия самой shell.

Дополнительно: реализовать парсинг && и || (условное выполнение команд), подстановку переменных окружения $VAR, поддержку редиректов >/< для вывода в файл и чтения из файла.

Основной упор необходимо делать на реализацию базового функционала (exec, builtins, pipelines). Проверять надо как интерактивно, так и скриптом. Код должен работать без ситуаций гонки, корректно освобождать ресурсы.

Совет: используйте пакеты os/exec, bufio (для ввода), strings.Fields (для разбиения командной строки на аргументы) и системные вызовы через syscall, если потребуется.*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	path, err := filepath.Abs("")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(path + ">")

	//при новом вводе определяется команда в соответствии с заданием
	for scanner.Scan() {
		input := scanner.Text()
		args := strings.Split(input, " ")
		switch args[0] {
		case "cd":
			err := os.Chdir(args[1])
			if err != nil {
				fmt.Println("неправильный путь")
			}
		case "pwd":
			fmt.Println(path)
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "kill":
			pid, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				break
			}
			prc, err := os.FindProcess(pid)
			if err != nil {
				fmt.Println(err)
				break
			}
			err = prc.Kill()
			if err != nil {
				fmt.Println(err)
				break
			}
		case "ps":
			prc, err := ps.Processes()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			procs := ""
			for _, value := range prc {
				procs += fmt.Sprintf("%s\n%d\n", value.Executable(), value.Pid())
			}
			fmt.Println(procs)
		default:
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		//в соответствии с заданной ранее командой происходит обновление дирректории
		path, err = filepath.Abs("")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(path + ">")
	}
}

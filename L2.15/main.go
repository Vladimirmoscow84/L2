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
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func main() {
	// сканер для чтения ввода пользователя с stdin
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// текущая рабочуая директория для вывода
		path, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(path + ">")

		// заверешение shell при команде от изера Ctrl+D
		if !scanner.Scan() {
			fmt.Println()
			os.Exit(0)
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		input = expandEnvVars(input)

		// контекст для возможности отмены команды (Ctrl+C)
		ctx, cancel := context.WithCancel(context.Background())

		// Канал для перехвата сигналов
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGQUIT) // Ctrl+C и Ctrl+\

		// в отдельной горутине слушаем канал для перехвата сгналов от юзера
		go func() {
			s := <-sigCh
			if s == syscall.SIGQUIT {
				os.Exit(0)
			}
			cancel()
		}()

		exitCode := runCommand(ctx, input)

		fmt.Printf("exit code: %d\n", exitCode)

		// закрытие канала и выход из контекста
		close(sigCh)
		cancel()
	}
}

// expandEnvVars  - подставляет переменные окружения $VAR в командной строке
func expandEnvVars(line string) string {
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, "$") && len(part) > 1 {
			val := os.Getenv(part[1:])
			line = strings.ReplaceAll(line, part, val)
		}
	}
	return line
}

// runCommand  - брабатывает команду, включая пайпы и встроенные команды
func runCommand(ctx context.Context, input string) int {
	if strings.Contains(input, "|") {
		return runPipeline(ctx, strings.Split(input, "|")) // если есть пайп, вызываем отдельную функцию
	}

	args := strings.Split(input, " ")

	// обработка встроенных команд
	switch args[0] {
	case "cd":
		return cmdCd(args)
	case "pwd":
		return cmdPwd(os.Stdout)
	case "echo":
		return cmdEcho(args, os.Stdout)
	case "kill":
		return cmdKill(ctx, args)
	case "ps":
		return cmdPs(ctx, args, os.Stdout)
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return 0
}

// cmdCd — встроенная команда cd
func cmdCd(args []string) int {
	if len(args) < 2 {
		fmt.Println("cd: missing operand")
		return 1
	}
	path := args[1]
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}
	if err := os.Chdir(path); err != nil {
		fmt.Println("cd:", err)
		return 1
	}
	return 0
}

// cmdPwd — встроенная команда pwd
func cmdPwd(w io.Writer) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, "pwd:", err)
		return 1
	}
	fmt.Fprintln(w, dir)
	return 0
}

// cmdEcho — встроенная команда echo
func cmdEcho(args []string, w io.Writer) int {
	fmt.Fprintln(w, strings.Join(args[1:], " "))
	return 0
}

// cmdKill — встроенная команда kill
func cmdKill(ctx context.Context, args []string) int {
	if len(args) < 2 {
		fmt.Println("kill: missing pid")
		return 1
	}

	cmdName := ""
	if runtime.GOOS == "windows" {
		cmdName = "taskkill"
		args = append([]string{"/PID"}, args[1:]...)
		args = append(args, "/F")
	} else {
		cmdName = "kill"
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// cmdPs — встроенная команда ps / tasklist
func cmdPs(ctx context.Context, args []string, w io.Writer) int {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "tasklist")
	} else {
		cmd = exec.CommandContext(ctx, "ps", args[1:]...)
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// runPipeline -  выполняет команды с пайпами
func runPipeline(ctx context.Context, cmds []string) int {
	var commands []*exec.Cmd
	for _, c := range cmds {
		args := strings.Fields(strings.TrimSpace(c))
		if len(args) == 0 {
			return 1
		}
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		commands = append(commands, cmd)
	}

	// cоединяем stdout каждой команды с stdin следующей
	for i := 0; i < len(commands)-1; i++ {
		stdout, _ := commands[i].StdoutPipe()
		commands[i+1].Stdin = stdout
	}

	// последняя команда выводит в stdout
	commands[len(commands)-1].Stdout = os.Stdout
	commands[len(commands)-1].Stderr = os.Stderr

	// запуск всех команд
	for _, cmd := range commands {
		start := time.Now()
		if err := cmd.Run(); err != nil {
			fmt.Println("error:", err)
			return 1
		}
		dur := time.Since(start)
		fmt.Printf("end cmd: %v, dur: %v\n", cmd, dur)
	}

	return 0
}

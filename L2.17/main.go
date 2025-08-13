/*
Утилита telnet (примитивный telnet-клиент)
Реализовать простой telnet-клиент с возможностью соединяться к TCP-серверу и взаимодействовать с ним:

Программа должна принимать параметры: хост, порт и опционально таймаут соединения (через флаг --timeout, по умолчанию 10 секунд).

После запуска, telnet-клиент устанавливает TCP-соединение с указанным host:port.

Все, что пользователь вводит в STDIN, должно отправляться в сокет; все, что приходит из сокета — печататься в STDOUT.

При нажатии комбинации клавиш Ctrl+D клиент должен закрыть соединение и завершиться. Если сервер закрыл соединение, клиент тоже завершается.

В случае, если попытка подключения не удалась (например, сервер недоступен) — программа завершается через заданный timeout с соответствующим сообщением об ошибке.

Проверить программу можно, например, подключившись к какому-нибудь публичному echo-серверу или SMTP (порт: 25) и вручную отправляя команды.

Обратите внимание на обработку буферов: желательно запускать чтение/запись в отдельных горутинах (для конкурентного ввода/вывода). Код должен быть без гонок. Реализация данной утилиты подразумевает использование пакета net (тип net.Conn), и возможно bufio для удобства чтения/записи.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type connection struct {
	host, port string
	timeout    time.Duration
}

// parseFflags - распарсивает аргументы командной строки
func parseFflags() (connection, error) {
	//переменные для хранения значений флага
	var h, p string
	var t time.Duration

	//флаги командной строки
	flag.DurationVar(&t, "timeout", 10*time.Second, "Таймаут соединения")
	flag.Parse()
	if len(flag.Args()) != 2 {
		return connection{}, fmt.Errorf("неверное количество аргументов командной строки")
	}
	h = flag.Arg(0)
	p = flag.Arg(1)

	args := connection{
		host:    h,
		port:    p,
		timeout: t,
	}
	return args, nil
}

// подключение к серверу
func (c *connection) client() error {
	con, err := net.DialTimeout("tcp", c.host+":"+c.port, c.timeout)
	if err != nil {
		return err
	}
	defer func() {
		//при выходе из программы закрываем соединение
		fmt.Println("Закрытие соединения")
		con.Close()
	}()

	//инрициализирование канала для сигналов ОС и отлова ошибок
	sysCh, errCh := make(chan os.Signal, 1), make(chan error)

	signal.Notify(sysCh, syscall.SIGINT)

	//чтение сообщений от хоста
	reader, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Println(reader)
	scanner := bufio.NewScanner(os.Stdin)

	//получение сообщений со стандартонго ввода и отправка их хосту, где они будут преобразованы
	go func() {
		for scanner.Scan() {
			_, err := con.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				errCh <- err
				close(errCh)
				break
			}
			l, err := bufio.NewReader(con).ReadString('\n')
			if err != nil {
				errCh <- err
				close(errCh)
				break
			}
			os.Stdout.Write([]byte(l))
		}
	}()

	select {
	case <-sysCh:
		fmt.Println("получен сигнал")
		return nil
	case err := <-errCh:
		return err

	}
	return nil
}

func main() {
	con, err := parseFflags()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = con.client()
	if err != nil {
		fmt.Println(err)
		return
	}

}

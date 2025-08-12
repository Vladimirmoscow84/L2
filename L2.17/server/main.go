package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// сервер по протоколу TCP на порту 8080
func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("===>Сервер запущен<===")

	// в бесконечном цикле сообщаем пользователю о подключении при получении соединения
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = connection.Write([]byte("Вы подключены к localhost:8080. Давайте разделим Ваш текст знаками '*'.\n"))
		if err != nil {
			log.Fatalln(err)
		}
		//сообщения от клиента разделяем знаками "*"
		for {
			str, err := bufio.NewReader(connection).ReadString('\n')
			//если клиент отключился, то закрывается соединение и ожидается соединение со следующим клиентом
			if err != nil {
				connection.Close()
				break
			}
			fmt.Println("Введеная пользователем строка:", str)

			modifyStr := strings.ReplaceAll(str, " ", "*")
			str = "Модифицированная строка: " + modifyStr
			_, err = connection.Write([]byte(str))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

/*Получение точного времени (NTP)
Создать программу, печатающую точное текущее время с использованием NTP-сервера.

Реализовать проект как модуль Go.

Использовать библиотеку ntp для получения времени.

Программа должна выводить текущее время, полученное через NTP (Network Time Protocol).

Необходимо обрабатывать ошибки библиотеки: в случае ошибки вывести её текст в STDERR и вернуть ненулевой код выхода.

Код должен проходить проверки (vet и golint), т.е. быть написан идиоматически корректно.*/

//воспользуемся документацией по адресу https://pkg.go.dev/github.com/beevik/ntp#section-readme

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime, err := getTime()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Точное текущее время: %v", currentTime)

}

func getTime() (time.Time, error) {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}
	return t, nil

}

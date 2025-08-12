/*Утилита wget (упрощенная)
Реализовать утилиту загрузки веб-страниц вместе со всем вложенным контентом (ресурсы, ссылки), аналогичную wget -m (мирроринг сайта).

Требования
Программа должна принимать URL и, возможно, глубину рекурсии (количество уровней ссылок, которые нужно скачать).

Должна уметь скачивать HTML-страницы, сохранять их локально, а также рекурсивно скачивать ресурсы: CSS, JS, изображения и т.д., а так же страницы, на которые есть ссылки (в рамках того же домена).

На выходе должен получиться локальный каталог, содержащий копию сайта (или его части), чтобы страницу можно было открыть офлайн.

Необходимо обрабатывать различные нюансы: относительные и абсолютные ссылки, дублирование (не скачивать один и тот же ресурс несколько раз), корректно формировать локальные пути для сохранения, избегать зацикливания по ссылкам.

Опционально: поддержать параллельное скачивание (например, ограничить до N одновременных загрузок), управлять robots.txt и пр.

Эта задача проверяет навыки сетевого программирования (HTTP-запросы), работы с файлами и строками, а также проектирования (нужно спланировать структуру, как хранить информацию о посещенных URL, как сохранять файлы и менять ссылки внутри HTML на локальные и т.д.).

Постарайтесь разбить программу на функции и пакеты: например, парсер HTML, загрузчик и т.п.

Обязательно учтите обработку ошибок (сетевых, файловых) и время выполнения (можно добавить таймауты на запросы).*/

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/opesun/goquery"
)

// downloadWebsite - скачивает статические ресурсы вэбсайта
func downloadWebsite(site string) error {

	//запрос на сайт
	response, err := http.Get(site)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//проверка успешности запроса по статускоду
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("неверный ответ сайта, код ответа: %v", response.StatusCode)
	}

	//запись тела запроса в html файл
	filename := strings.Split(site, "/")[2] + ".html"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// downloadResources -  загружает изображения и описания внешнего вида сайта
func downloadResources(fileName string, url string) error {
	// получение данных по URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Создание файлф с переданным именем и запись в него полученных данных
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// parseWebsite - парсит ресурсы вэбсайта
func parseWebsite(site string) error {
	data, err := goquery.ParseUrl(site)
	if err != nil {
		return err
	}
	for _, url := range data.Find("").Attrs("href") {
		var str []string
		switch {
		case strings.Contains(url, ".png"):
			str = strings.Split(url, "/")
			downloadResources(str[len(str)-1], url)
		case strings.Contains(url, ".jpg"):
			str = strings.Split(url, "/")
			downloadResources(str[len(str)-1], url)
		case strings.Contains(url, ".css"):
			str = strings.Split(url, "/")
			downloadResources(str[len(str)-1], url)
		}
	}
	return nil
}

func main() {
	url := flag.String("s", "https://gazeta.ru", "site url")
	flag.Parse()
	if strings.Contains(*url, "https://") || strings.Contains(*url, "http://") {
		if err := downloadWebsite(*url); err != nil {
			fmt.Println(err)
			return
		}
		if err := parseWebsite(*url); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("invalid url")
		return
	}
}

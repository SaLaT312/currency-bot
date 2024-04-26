package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/robfig/cron"
	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	XMLName   xml.Name `xml:"Valute"`
	ID        string   `xml:"ID,attr"`
	NumCode   string   `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   string   `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     string   `xml:"Value"`
	VunitRate string   `xml:"VunitRate"`
}

func GetDataFromBank() {

	currentTime := time.Now()
	formattedDate := currentTime.Format("02/01/2006")

	RedisCheck() // Проверяем подключение к Redis

	url := "https://www.cbr.ru/scripts/XML_daily.asp?date_req=" + formattedDate // URL ЦБ РФ с датой запроса

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании GET-запроса:", err)
		return
	}

	req.Header = http.Header{
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:123.0) Gecko/20100101 Firefox/123.0"}, // Необходимый доп заголовок для запроса
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Ошибка при выполнении GET-запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return
	}

	r := bytes.NewReader(body)
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs

	err = d.Decode(&valCurs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	log.Printf("Actual Data Date: %s\n", valCurs.Date)

	RedisSet("Date", valCurs.Date)

	ListOfValuteCode := make([]string, 0)

	for _, valute := range valCurs.Valutes {

		ListOfValuteCode = append(ListOfValuteCode, valute.CharCode)

		RedisSet(valute.CharCode+"_Code", valute.CharCode)                              // Код валюты
		RedisSet(valute.CharCode+"_Name", valute.Name)                                  // Полное имя валюты
		RedisSet(valute.CharCode+"_Nominal", valute.Nominal)                            // Коэф для преобразования
		RedisSet(valute.CharCode+"_Value", strings.Replace(valute.Value, ",", ".", -1)) // Курс валюты в формате float

	}
	ListOfValuteCodeForRedis := strings.Join(ListOfValuteCode, ",") // Подготовка списка полученный валют и запись его в редис
	//log.Println("Write ListOfValuteCodeForRedis to Redis:")
	//log.Println(ListOfValuteCodeForRedis)
	RedisSet("ListOfValuteCode", ListOfValuteCodeForRedis)

}

func CronTabFunc() {

	c := cron.New()
	c.AddFunc("@daily", func() { GetDataFromBank() })
	c.Start()

}

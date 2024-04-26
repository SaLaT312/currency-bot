package main

import (
	"fmt"
	"strconv"
	"strings"
)

func CheckTypeFunc(checkValue string) (typeValue string) {

	accessibleValutes := strings.Split(RedisGet("ListOfValuteCode"), ",")

	// Проверка наличия строки в массиве валют
	for _, str := range accessibleValutes {
		if str == strings.ToUpper(checkValue) {
			typeValue = "valute" // Если нашли, возвращаем формат валюты
			return typeValue
		}
	}

	// Попытка преобразовать введенное значение в float64
	out := strings.Replace(checkValue, ",", ".", -1)

	_, err := strconv.ParseFloat(out, 64)
	if err != nil {
		typeValue = "error" // Если не получилось, возвращаем неизвестный формат
		return typeValue
	}

	typeValue = "number" // Если получилось, возвращаем формат числа
	return typeValue
}

func Actual() (textMessage string) {

	ActualDate := RedisGet("Date")
	textMessage = fmt.Sprintf("Actual date: %s", ActualDate)

	return
}

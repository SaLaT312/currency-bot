package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func SetValueForUser(UserID int64, ImportText string) (textMessage string) {

	ConvertValute := strings.ToUpper(ImportText)

	RedisSet(strconv.Itoa(int(UserID)), ConvertValute)

	textMessage = "Валюта для конвертации - " + ConvertValute

	return
}

func GetValueForUser(UserID int64) (textMessage string) {

	ConvertValute := RedisGet(strconv.Itoa(int(UserID)))

	if ConvertValute != "" {
		textMessage = ConvertValute
	} else {
		textMessage = "USD"
	}

	return
}

func GetCurrencies() (textMessage string) {

	ListOfValute := strings.Split(RedisGet("ListOfValuteCode"), ",")

	for _, ValuteCode := range ListOfValute {

		ValuteName := RedisGet(ValuteCode + "_Name")

		textMessage = textMessage + fmt.Sprintf("%s - %s\n", ValuteCode, ValuteName)
	}

	return
}

func GetRates() (textMessage string) {

	ListOfValute := strings.Split(RedisGet("ListOfValuteCode"), ",")

	for _, ValuteCode := range ListOfValute {

		ValuteName := RedisGet(ValuteCode + "_Name")
		ValuteNominal := RedisGet((ValuteCode + "_Nominal"))
		ValuteValue, _ := strconv.ParseFloat(RedisGet(ValuteCode+"_Value"), 8)

		textMessage = textMessage + fmt.Sprintf("%s %s (%s) = %.2f RUB\n", ValuteNominal, ValuteCode, ValuteName, ValuteValue)
	}

	return
}

func round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func ConvertValute(UserID int64, ImportText string) (textMessage string) {

	UserValute := GetValueForUser(UserID)

	ImportTextReady := strings.Replace(ImportText, ",", ".", -1)
	ValuteData, _ := strconv.ParseFloat(ImportTextReady, 8)

	ValuteNominal := RedisGet((UserValute + "_Nominal"))
	ValuteValue := RedisGet((UserValute + "_Value"))

	value := round(ValuteData, 2)

	cost, _ := strconv.ParseFloat(ValuteValue, 8)
	nominal, _ := strconv.ParseFloat(ValuteNominal, 8)

	result := (value * cost) / nominal
	formattedResult := formatFloatToString(result, " ", 2)

	ImportText = strings.Replace(ImportText, ".", ",", -1)

	textMessage = fmt.Sprintf("%s %s = %s RUB", ImportText, UserValute, formattedResult)

	return
}

// Функция для форматирования числа типа float64 в строку с заданным разделителем тысяч и заданным количеством знаков после запятой
func formatFloatToString(num float64, thousandsSeparator string, precision int) string {
	// Преобразуем число в строку с заданным количеством знаков после запятой
	str := strconv.FormatFloat(num, 'f', precision, 64)

	// Разделяем строку на целую и дробную части
	parts := splitDecimal(str)

	// Добавляем разделитель тысяч
	formattedInteger := addThousandSeparator(parts[0], thousandsSeparator)

	// Объединяем целую и дробную части с разделителем дробной части
	formattedNum := formattedInteger + "," + parts[1]

	return formattedNum
}

// Функция для разделения числа с плавающей точкой на целую и дробную части
func splitDecimal(numStr string) []string {
	parts := make([]string, 2)
	dotIndex := -1
	for i, char := range numStr {
		if char == '.' {
			dotIndex = i
			break
		}
	}
	if dotIndex == -1 {
		parts[0] = numStr
		parts[1] = "0"
	} else {
		parts[0] = numStr[:dotIndex]
		parts[1] = numStr[dotIndex+1:]
	}
	return parts
}

// Функция для добавления разделителя тысяч в целую часть числа
func addThousandSeparator(integerPart, separator string) string {
	n := len(integerPart)
	if n <= 3 {
		return integerPart
	}
	return addThousandSeparator(integerPart[:n-3], separator) + separator + integerPart[n-3:]
}

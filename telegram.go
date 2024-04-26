package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TelegramBot(TelegramApiToken string) {
	bot, err := tgbotapi.NewBotAPI(TelegramApiToken)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var textMessage string

	for update := range updates {
		// Пропускаем неподдерживаемые типы обновлений
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {

		case "/start":
			textMessage = "Добро пожаловать в бот для конвертации валют. Для получения списка команд воспользуйтесь справкой /help"

		case "/help":
			textMessage = HelpText() // Вывод текст помощи

		case "/get_currencies":
			textMessage = GetCurrencies()

		case "/get_rates":
			textMessage = GetRates() //

		case "/actual":
			textMessage = Actual() //

		case "/get_def_rate":
			textMessage = "Валюта для конвертации - " + GetValueForUser(update.Message.From.ID)

		default:
			switch CheckTypeFunc(update.Message.Text) { // Отдать в обработку анализатору сообщений

			case "valute":
				textMessage = SetValueForUser(update.Message.From.ID, update.Message.Text)
				//
			case "number":
				textMessage = ConvertValute(update.Message.From.ID, update.Message.Text)
				//
			default:
				textMessage = fmt.Sprintf("Введено неверное значение: %s\nДля справки воспользуйтесь командой /help", update.Message.Text)
				//

			}

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
		_, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}

	}
}

func HelpText() string {

	return `Available commands:
/help - Displays this help
/get_currencies - Display available currencies
/get_rates - Request the rate of all currencies
/get_def_rate - Get the current currency
/actual - Check actual date of currency

To convert, enter a numeric value (Use a dot or comma to separate the integer and fractional parts of the number).
To change the conversion currency, enter its code. Currency codes can be obtained using the /get_rates command. Input is available in lowercase and uppercase letters.
The current currency for conversion can be found by the command /get_def_rate
------------------
Доступны команды:
/help - Выводит эту справку
/get_currencies - Вывести доступные валюты
/get_rates - Запрос курса всех валют
/get_def_rate - Узнать текущую валюту
/actual - показывает дату актуальности валют

Для конвертации введите числовое значение (Используйте точку или запятую для разделения целой и дробной части числа).
Для смены валюты конвертации введите ее код. Коды валют можно получить по команде /get_rates Доступен ввод строчным и прописными буквами.
Текущую валюту для конвертации можно узнать по команде /get_def_rate`

}

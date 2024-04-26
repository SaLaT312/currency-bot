package main

import (
	"log"

	"github.com/go-redis/redis"
)

// Создаем подключение к Redis
var RedisAddr, RedisPass string
var RedisDB int

var ClientRedis = redis.NewClient(&redis.Options{
	Addr:     RedisAddr,
	Password: RedisPass, // пустой пароль
	DB:       RedisDB,   // база данных 0
})

func RedisCheck() {

	// Проверяем подключение к Redis
	pong, err := ClientRedis.Ping().Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	log.Printf("Подключение к Redis успешно: %s", pong)

}

func RedisSet(KeyRedis, DataRedis string) {

	err := ClientRedis.Set(KeyRedis, DataRedis, 0).Err()
	if err != nil {
		log.Fatalf("Ошибка записи ключа %s в Redis: %v", KeyRedis, err)
	}

}

func RedisGet(KeyRedis string) (KeyData string) {

	KeyData, _ = ClientRedis.Get(KeyRedis).Result()
	//if KeyData == "" {
	//	log.Printf("Ошибка получения ключа %s из Redis: %v", KeyRedis, KeyData)
	//}

	return KeyData
}

package main

import (
	"flag"
	"iRo/internal/application"
	"log"
)

// значения меняем при сборке, см. Makefile
var (
	version = "0.0.0"
	commit  = "unset"
)

// константы
var (
	sysCfgDefault = "./config/system.yaml"
)

// main - запуск приложения
func main() {
	// сообщаем о запуске
	log.Print("запуск версии <", version, "V", commit, ">\n")
	// расположение файла настройки
	sysCfgPath := flag.String("config", sysCfgDefault, "путь доступа и имя файла конфигурации системы")
	sysCfgShow := flag.Bool("info", false, "выдача настроек")
	// парсинг аргументов командной строки
	flag.Parse()
	// в случае ошибки
	if err := application.Run(*sysCfgShow, sysCfgPath); err != nil {
		// показываем сообщение и завершаем программу
		log.Print("остановлено после сбоя <", err, ">\n")
	}
	// сообщаем о завершении работы
	log.Print("завершено\n")
}

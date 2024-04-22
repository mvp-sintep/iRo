package main

import (
	"flag"
	"fmt"
	"iRo/internal/application"
	"log"
	"time"
)

// значения меняем при сборке, см. Makefile
var (
	version = "0.0.0"
	commit  = "unset"
)

// константы
var (
	sysCfgDefault = "./config/system.yaml"
	aeCfgDefault  = "./config/ae.yaml"
)

// main - запуск приложения
func main() {
	// сообщаем о запуске
	log.Print("запуск версии <", version, "V", commit, ">\n")
	// расположение файла настройки
	sysCfgPath := flag.String("config", sysCfgDefault, "путь доступа и имя файла конфигурации системы")
	aeCfgPath := flag.String("ae", aeCfgDefault, "путь доступа и имя файла конфигурации модуля тревог и событий")
	cfgShow := flag.Bool("info", false, "выдача настроек")
	// парсинг аргументов командной строки
	flag.Parse()
	// сообщаем о начале работы
	fmt.Print(time.Now().Format("2006/01/02 15:04:05"), " ожидается команда, нажмите Ctrl^C для выхода...\n")
	// в случае ошибки
	if err := application.Run(*cfgShow, sysCfgPath, aeCfgPath); err != nil {
		// показываем сообщение и завершаем программу
		log.Print("остановлено после сбоя <", err, ">\n")
	}
	// сообщаем о завершении работы
	fmt.Print("\n", time.Now().Format("2006/01/02 15:04:05"), " выход\n")
}

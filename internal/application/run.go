package application

import (
	"fmt"
	"iRo/internal/config"

	"gopkg.in/yaml.v3"
)

// Run - запуск приложения
func Run(cfgShow bool, sysCfgPath *string) error {
	// создаем запись конфигурации системы
	sysCfg := config.New()
	// загружаем данные
	if err := sysCfg.Load(sysCfgPath); err != nil {
		return err
	}
	// выдача настроек
	if cfgShow {
		out, err := yaml.Marshal(sysCfg)
		if err != nil {
			return err
		}
		fmt.Printf("настройка:\n%s", string(out))
	}
	// нет ошибок
	return nil
}

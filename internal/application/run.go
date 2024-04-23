package application

import (
	"context"
	"fmt"
	"iRo/internal/config"
	"iRo/internal/daemon"

	"gopkg.in/yaml.v3"
)

// Run - запуск приложения
func Run(cfgShow bool, sysCfgPath *string, aeCfgPath *string) error {
	// создаем запись конфигурации системы
	sysCfg := config.NewSystemConfig()
	// загружаем данные
	if err := sysCfg.Load(sysCfgPath); err != nil {
		return err
	}
	// создаем запись конфигурации модуля тревог и событий
	aeCfg := config.NewAEConfig()
	// загружаем данные
	if err := aeCfg.Load(aeCfgPath); err != nil {
		// выход при ошибке
		return err
	}
	// выдача настроек
	if cfgShow {
		out, err := yaml.Marshal(sysCfg)
		if err != nil {
			return err
		}
		fmt.Printf("настройка:\n%s", string(out))
		out, err = yaml.Marshal(aeCfg)
		if err != nil {
			return err
		}
		fmt.Printf("%s", string(out))
	}
	// создаем фоновый процесс
	daemon, err := daemon.New(context.Background(), sysCfg, aeCfg)
	// проверяем
	if err != nil {
		return err
	}
	// запускаем и возвращаем ошибку
	if err := daemon.Run(); err != nil {
		return err
	}
	// нет ошибок
	return nil
}

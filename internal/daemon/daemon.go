package daemon

import (
	"context"
	"iRo/internal/config"
	"iRo/internal/core"
	"iRo/internal/driver/com"
	"iRo/internal/server/mbrtu"
	"iRo/internal/server/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Daemon - фоновый процесс верхнего уровня
type Daemon struct {
	context  context.Context
	cancel   context.CancelFunc
	core     []byte
	comDrv   *com.Driver
	mbrtuSrv *mbrtu.Server
	webECh   chan error
	webSrv   *web.Server
}

// New - создание процесса
func New(ctx context.Context, sysCfg *config.SystemConfiguration) (*Daemon, error) {
	var err error
	// создаем процесс
	daemon := &Daemon{
		webECh: make(chan error),
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	daemon.context, daemon.cancel = context.WithCancel(ctx)
	// ядро данных
	daemon.core = core.Data[:]

	// создаем драйвер COM порта
	if daemon.comDrv, err = com.New(daemon.context, &sysCfg.COM[0]); err != nil {
		return nil, err
	}
	// создаем сервер MODBUS RTU
	if daemon.mbrtuSrv, err = mbrtu.New(daemon.context, &daemon.core, daemon.comDrv, &sysCfg.Modbus.RTU[0]); err != nil {
		return nil, err
	}
	// создаем HTTP сервер
	if daemon.webSrv, err = web.New(daemon.context, daemon.webECh, &daemon.core, &sysCfg.HTTP); err != nil {
		return nil, err
	}
	// вернем указатель на процесс
	return daemon, nil
}

// Run - запуск процесса
func (o *Daemon) Run() error {
	// чистим ядро
	for i := range o.core {
		o.core[i] = 0
	}
	// запускаем сервер MODBUS RTU до запуска драйвера (иначе возможна блокировка из-за сигнала DSR)
	go func() { o.mbrtuSrv.Run() }()
	// запускаем HTTP сервер
	go func() { o.webECh <- o.webSrv.Run() }()
	// запускаем драйвер COM порта
	go func() { o.comDrv.Run() }()
	// создаем канал для системных сигналов
	sch := make(chan os.Signal, 1)
	// подписываемся на сигналы CTRL^C и Ctrl^J
	signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM)
	// ждем сигнал
	select {
	// системный канал
	case <-sch:
		// завершаем работу
		o.Shutdown()
		// выход без ошибки
		return nil
	// канал ошибки HTTP сервера
	case err := <-o.webECh:
		// показываем сообщение
		log.Print("ошибка HTTP сервера <", err, ">")
		// завершаем работу
		o.Shutdown()
		// возвращаем ошибку
		return err
	}
}

// Shutdown - завершение процесса
func (o *Daemon) Shutdown() {
	// завершение работы HTTP сервера
	if err := o.webSrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка HTTP сервера <", err, ">")
	}
	// завершение работы сервера MODBUS RTU
	if err := o.mbrtuSrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка сервера MODBUS RTU <", err, ">")
	}
	// завершение работы драйвера RTU
	if err := o.comDrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка драйвера COM порта <", err, ">")
	}
	// после выхода дадим команду на закрытие контекста
	defer o.cancel()
}

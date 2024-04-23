package daemon

import (
	"context"
	"iRo/internal/config"
	"iRo/internal/core"
	"iRo/internal/daemon/ae"
	"iRo/internal/driver/com"
	"iRo/internal/driver/db"
	"iRo/internal/driver/tcp"
	"iRo/internal/server/mbrtu"
	"iRo/internal/server/mbtcp"
	"iRo/internal/server/ua"
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
	nda      chan struct{}
	comDrv   *com.Driver
	tcpDrv   *tcp.Driver
	tcpECh   chan error
	mbrtuSrv *mbrtu.Server
	mbtcpSrv *mbtcp.Server
	webECh   chan error
	webSrv   *web.Server
	uaSrv    *ua.Server
	dbPool   *db.Pool
	ae       *ae.Service
}

// New - создание процесса
func New(ctx context.Context, sysCfg *config.SystemConfiguration, aeCfg *config.AEConfig) (*Daemon, error) {
	var err error
	// создаем процесс
	daemon := &Daemon{
		nda:    make(chan struct{}),
		tcpECh: make(chan error),
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
	// контроль изменения данных ядра
	daemon.nda = 0

	// создаем драйвер COM порта
	if daemon.comDrv, err = com.New(daemon.context, &sysCfg.COM[0]); err != nil {
		return nil, err
	}
	// создаем драйвер TCP порта
	if daemon.tcpDrv, err = tcp.New(daemon.context, daemon.tcpECh, &sysCfg.Modbus.TCP); err != nil {
		return nil, err
	}
	// создаем сервер MODBUS RTU
	if daemon.mbrtuSrv, err = mbrtu.New(daemon.context, &daemon.core, daemon.nda, daemon.comDrv, &sysCfg.Modbus.RTU[0]); err != nil {
		return nil, err
	}
	// создаем сервер MODBUS TCP
	if daemon.mbtcpSrv, err = mbtcp.New(daemon.context, &daemon.core, daemon.nda, daemon.tcpDrv, &sysCfg.Modbus.TCP); err != nil {
		return nil, err
	}
	// создаем пул соединений с СУБД
	if daemon.dbPool, err = db.New(daemon.context, &sysCfg.DB); err != nil {
		return nil, err
	}
	// создаем HTTP сервер
	if daemon.webSrv, err = web.New(daemon.context, daemon.webECh, &daemon.core, daemon.dbPool, &sysCfg.HTTP); err != nil {
		return nil, err
	}
	// создаем UA сервер
	if daemon.uaSrv, err = ua.New(daemon.context, &sysCfg.UA); err != nil {
		return nil, err
	}
	// создаем сервис тревог и событий
	if daemon.ae, err = ae.New(daemon.context, &daemon.core, daemon.nda, daemon.dbPool, aeCfg); err != nil {
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
	// запускаем сервер MODBUS TCP до запуска драйвера (иначе возможна блокировка из-за сигнала DSR)
	go func() { o.mbtcpSrv.Run() }()
	// запускаем HTTP сервер
	go func() { o.webECh <- o.webSrv.Run() }()
	// запускаем UA сервер
	go func() { o.uaSrv.Run() }()
	// запускаем драйвер COM порта
	go func() { o.comDrv.Run() }()
	// запускаем драйвер TCP порта
	go func() { o.tcpDrv.Run() }()
	// запускаем сервис тревог и событий
	go func() { o.ae.Run() }()

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
	// канал ошибки TCP сервера
	case err := <-o.tcpECh:
		// показываем сообщение
		log.Print("ошибка TCP сервера <", err, ">")
		// завершаем работу
		o.Shutdown()
		// возвращаем ошибку
		return err
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
	// завершение работы UA сервера
	if err := o.uaSrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка OPC UA сервера <", err, ">")
	}
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
	// завершение работы сервера MODBUS TCP
	if err := o.mbtcpSrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка сервера MODBUS TCP <", err, ">")
	}
	// завершение работы драйвера TCP
	if err := o.tcpDrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка драйвера TCP порта <", err, ">")
	}
	// завершение работы драйвера RTU
	if err := o.comDrv.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка драйвера COM порта <", err, ">")
	}
	// завершение работы пула СУБД
	if err := o.dbPool.Shutdown(); err != nil {
		// показываем сообщение
		log.Print("ошибка пула СУБД <", err, ">")
	}
	// после выхода дадим команду на закрытие контекста
	defer o.cancel()
}

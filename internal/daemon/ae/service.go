package ae

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"
	"iRo/internal/driver/db"
	"log"
	"strconv"
	"strings"
	"time"
)

// Перечень состояний в соответствии с ISA 18.2
const (
	eventStateNORM  = 1 // Тревога не активна, подтверждена
	eventStateUNACK = 2 // Тревога активна, не подтверждена
	eventStateACKED = 3 // Тревога активна, подтверждена
	eventStateRTNUN = 4 // Тревога не активна, не подтверждена
	eventStateSHLVD = 5 // Активация тревоги отложена оператором
	eventStateDSUPR = 6 // Активация тревоги запрещена программой
	eventStateOOSRV = 7 // Тревога выведена из эксплуатации
	eventStateERR   = 8 // Ошибка обработки
)

// address - полный адрес события в базе и в ядре данных
type address struct {
	msg    string // идентификатор из таблицы msg
	offset int    // смещение в ядре данных
	mask   uint8  // битовая маска из поля byte mask
}

// event - источник события
type event struct {
	enable bool
	state  byte // код состояния , см. eventState... const
	active address
}

// Service - блок данных сервиса
type Service struct {
	context context.Context    // контекст для сервиса
	cancel  context.CancelFunc // функция закрытия контекста сервиса
	cfg     *config.AEConfig   // запись конфигурации
	core    *[]byte            // ядро данных
	nda     <-chan struct{}    // ядро данных изменено
	pool    *db.Pool           // СУБД
	Event   []event            // массив событий
}

// New - создание блока данных сервиса
func New(ctx context.Context, core *[]byte, nda <-chan struct{}, pool *db.Pool, cfg *config.AEConfig) (*Service, error) {
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации AE сервиса")
	}
	// создаем блок данных
	service := &Service{
		cfg:   cfg,
		core:  core,
		nda:   nda,
		pool:  pool,
		Event: make([]event, len(cfg.Event)),
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	service.context, service.cancel = context.WithCancel(ctx)
	// сервис создан
	return service, nil
}

// Run - запуск сервиса тревог и событий
func (o *Service) Run() {

	getA := func(cfg config.EventConfig) (offset int, mask uint8, err error) {
		str := strings.SplitN(cfg.Active, ":", 2)
		o, err := strconv.Atoi(str[0])
		if err != nil {
			return 0, 0, err
		}
		x, err := strconv.Atoi(str[1])
		if err != nil || x < 0 || x > 7 {
			return o, 0, err
		}
		return o, uint8(1) << x, nil
	}

	for i := 0; i < len(o.Event); i++ {
		o.Event[i].enable = false
		o.Event[i].state = eventStateNORM
		rows, err := o.pool.Query("select id from msg where descriptor -> 'text' ? '" + o.cfg.Event[i].Text + "'")
		if err != nil {
			rows.Close()
			continue
		}
		for rows.Next() {
			if vs, err := rows.Values(); err == nil && len(vs) == 1 {
				o.Event[i].active.msg = fmt.Sprintf("%d", vs[0].(int64))
				if o.Event[i].active.offset, o.Event[i].active.mask, err = getA(o.cfg.Event[i]); err != nil {
					continue
				}
				o.Event[i].enable = true
				break
			}
		}
		rows.Close()
	}

	if len(o.Event) == 0 {
		return
	}

	max := len(*o.core) - 1

	ns := false

	for {
		// signal selector
		select {

		// application terminate signal
		case <-o.context.Done():
			return

		// nda signal (new data avaiable)
		case <-o.nda:
			ns = true

		// search event if nda set true
		case <-time.After(250 * time.Millisecond):
			flag := false
			tmp := ""
			if ns {
				ns = false
				for i := 0; i < len(o.Event); i++ {
					if !o.Event[i].enable || o.Event[i].active.offset < 0 || o.Event[i].active.offset > max {
						continue
					}
					switch o.Event[i].state {
					case eventStateNORM:
						if ((*o.core)[o.Event[i].active.offset] & (*o.core)[o.Event[i].active.mask]) != 0 {
							o.Event[i].state = eventStateUNACK
							if flag {
								tmp += ","
							}
							flag = true
							tmp += fmt.Sprintf("(%v,%d,'%v')", o.Event[i].active.msg, o.Event[i].state, o.cfg.Event[i].Text)
						}
					case eventStateUNACK:
						if ((*o.core)[o.Event[i].active.offset] & (*o.core)[o.Event[i].active.mask]) == 0 {
							o.Event[i].state = eventStateNORM
							if flag {
								tmp += ","
							}
							flag = true
							tmp += fmt.Sprintf("(%v,%d,'%v')", o.Event[i].active.msg, o.Event[i].state, o.cfg.Event[i].Text)
						}
					default:
						o.Event[i].state = eventStateNORM
					}
				}
			}
			if flag {
				go func(str string) {
					rows, err := o.pool.Query("insert into event(msg,state,text) values " + str)
					if err != nil {
						log.Print(err)
					}
					rows.Close()
				}(tmp)
			}
		}
	}
}

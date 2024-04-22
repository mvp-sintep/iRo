package ae

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"
	"iRo/internal/core"
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
	pool    *db.Pool           // СУБД
	Event   []event            // массив событий
	nda     *int               // переменная для контроля обновления ядра данных
}

// New - создание блока данных сервиса
func New(ctx context.Context, pool *db.Pool, nda *int, cfg *config.AEConfig) (*Service, error) {

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

	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации AE сервиса")
	}

	// создаем блок данных
	service := &Service{
		cfg:   cfg,
		pool:  pool,
		Event: make([]event, len(cfg.Event)),
		nda:   nda,
	}

	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	service.context, service.cancel = context.WithCancel(ctx)

	for i := 0; i < len(service.Event); i++ {
		service.Event[i].enable = false
		service.Event[i].state = eventStateNORM
		rows, err := pool.Query("select id from msg where descriptor -> 'text' ? '" + cfg.Event[i].Text + "'")
		if err != nil {
			rows.Close()
			continue
		}
		for rows.Next() {
			if vs, err := rows.Values(); err == nil && len(vs) == 1 {
				service.Event[i].active.msg = fmt.Sprintf("%d", vs[0].(int64))
				if service.Event[i].active.offset, service.Event[i].active.mask, err = getA(cfg.Event[i]); err != nil {
					continue
				}
				service.Event[i].enable = true
				break
			}
		}
		rows.Close()
	}

	return service, nil
}

func (o *Service) Run() {

	if len(o.Event) == 0 {
		return
	}

	max := len(core.Data) - 1

	previousSecond := 58

	for {
		// signal selector
		select {
		// application terminate signal
		case <-o.context.Done():
			return
		// clear store
		case <-time.After(1 * time.Second):
			h, m, s := time.Now().Clock()
			if h == 11 && m == 59 && s == 59 && previousSecond == 58 {
				for i := 0; i < len(o.Event) && i < len(o.cfg.Event); i++ {
					if !o.Event[i].enable || o.cfg.Event[i].Store <= 0 {
						continue
					}
					rows, err := o.pool.Query(fmt.Sprintf("delete from event where (msg = %s and ts < now() - interval '%d day')", o.Event[i].active.msg, 10*o.cfg.Event[i].Store))
					if err != nil {
						log.Print(err)
					}
					rows.Close()
				}
			}
			previousSecond = s
		// search event if nda set true
		default:
			flag := false
			str := ""
			if *o.nda > 0 {
				*o.nda = 0
				for i := 0; i < len(o.Event); i++ {
					if !o.Event[i].enable || o.Event[i].active.offset < 0 || o.Event[i].active.offset > max {
						continue
					}
					switch o.Event[i].state {
					case eventStateNORM:
						if (core.Data[o.Event[i].active.offset] & core.Data[o.Event[i].active.mask]) != 0 {
							o.Event[i].state = eventStateUNACK
							if flag {
								str += ","
							}
							flag = true
							str += fmt.Sprintf("(%v,%d,'%v')", o.Event[i].active.msg, o.Event[i].state, o.cfg.Event[i].Text)
						}
					case eventStateUNACK:
						if (core.Data[o.Event[i].active.offset] & core.Data[o.Event[i].active.mask]) == 0 {
							o.Event[i].state = eventStateNORM
							if flag {
								str += ","
							}
							flag = true
							str += fmt.Sprintf("(%v,%d,'%v')", o.Event[i].active.msg, o.Event[i].state, o.cfg.Event[i].Text)
						}
					default:
						o.Event[i].state = eventStateNORM
					}
				}
			}
			if flag {
				rows, err := o.pool.Query("insert into event(msg,state,text) values " + str)
				if err != nil {
					log.Print(err)
				}
				rows.Close()
			}
		}
	}
}

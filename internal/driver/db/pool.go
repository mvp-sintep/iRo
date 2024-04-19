package db

import (
	"context"
	"errors"
	"fmt"
	"iRo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool - блок данных соединения с СУБД
type Pool struct {
	context context.Context    // контекст для сервера
	cancel  context.CancelFunc // функция закрытия контекста сервера
	cfg     *config.DBConfig   // запись конфигурации
	pgx     *pgxpool.Pool      // пул
}

// New - создание пула соединений с СУБД
func New(ctx context.Context, cfg *config.DBConfig) (*Pool, error) {
	var err error
	// проверяем аргументы
	if cfg == nil {
		return nil, errors.New("нет данных конфигурации СУБД")
	}
	// создаем блок данных
	pool := &Pool{
		cfg: cfg,
	}
	// если не задан контекст
	if ctx == nil {
		ctx = context.Background()
	}
	// задаем значения
	pool.context, pool.cancel = context.WithCancel(ctx)
	// создаем пул соединений
	pool.pgx, err = pgxpool.New(pool.context, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", pool.cfg.User, pool.cfg.Password, pool.cfg.Address, pool.cfg.Port, pool.cfg.Base))
	// проверяем на ошибку
	if err != nil {
		return nil, err
	}
	// проверяем на работоспособность соединения
	if _, err := pool.pgx.Exec(ctx, "SELECT 1;"); err != nil {
		return nil, err
	}
	// вернем пул
	return pool, nil
}

// Shutdown - завершение работы
func (o *Pool) Shutdown() error {
	// закроем контекст пула
	defer o.cancel()
	// выход без ошибки
	return nil
}

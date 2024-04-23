package config

import (
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

// DataConfig - типовой блок конфигурации области чтения/записи данных
type DataConfig struct {
	Start int `yaml:"start"` // смещение начала блока данных
	Bytes int `yaml:"bytes"` // кол-во байт данных
}

// COMPortConfig - данные конфигурации COM порта
type COMPortConfig struct {
	File     string `yaml:"file"` // имя пвсевдо файла устройства
	BaudRate int    `yaml:"baudrate"`
	DataBits int    `yaml:"databits"`
	Parity   string `yaml:"parity"`
	StopBits int    `yaml:"stopbits"`
	Timeout  int    `yaml:"timeout"`
}

// ModbusRTUConfig - данные конфигурации MODBUS RTU устройства
type ModbusRTUConfig struct {
	Node int        `yaml:"node"`
	Core DataConfig `yaml:"core"`
}

// ModbusTCPConfig - данные конфигурации MODBUS TCP порта
type ModbusTCPConfig struct {
	Address string     `yaml:"address"`
	Port    int        `yaml:"port"`
	Read    int        `yaml:"read"`
	Write   int        `yaml:"write"`
	Control int        `yaml:"control"`
	Core    DataConfig `yaml:"core"`
}

// ModbusConfig - каталог записей конфигурации MODBUS
type ModbusConfig struct {
	RTU []ModbusRTUConfig `yaml:"rtu"`
	TCP ModbusTCPConfig   `yaml:"tcp"`
}

// UAConfig - настройки сервера opc ua
type UAConfig struct {
	Port      uint `yaml:"port"`
	Namespace int  `yaml:"namespace"`
}

// DBConfig - настройки сервера postgres
type DBConfig struct {
	Address  string `yaml:"address"`
	Proto    string `yaml:"proto"`
	Port     string `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
	Base     string `yaml:"base"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// HTTPConfig - настройки сервера http
type HTTPConfig struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	Read     int    `yaml:"read"`
	Write    int    `yaml:"write"`
	Shutdown int    `yaml:"shutdown"`
	Root     string `yaml:"root"`
}

// SystemConfig - все данные настройки
type SystemConfig struct {
	COM    []COMPortConfig `yaml:"com"`
	Modbus ModbusConfig    `yaml:"modbus"`
	UA     UAConfig        `yaml:"ua"`
	DB     DBConfig        `yaml:"postgres"`
	HTTP   HTTPConfig      `yaml:"http"`
}

// NewSystemConfig - создание записи конфигурации
func NewSystemConfig() *SystemConfiguration {

	address := func() string { // получим первый IP адрес первой сетевой платы
		if list, err := net.InterfaceAddrs(); err == nil {
			for _, x := range list {
				if i, ok := x.(*net.IPNet); ok && !i.IP.IsLoopback() {
					if i.IP.To4() != nil {
						return i.IP.String()
					}
				}
			}
		}
		return "localhost"
	}()

	// создадим запись конфигурации со значениями по умолчанию
	// для настроек доступа к postgres предусмотрено использование переменных окружения
	return &SystemConfig{
		COM: []COMPortConfig{
			{
				File:     "/dev/ttyS0",
				BaudRate: 115200,
				DataBits: 8,
				Parity:   "N",
				StopBits: 1,
				Timeout:  1,
			},
		},
		Modbus: ModbusConfig{
			RTU: []ModbusRTUConfig{
				{
					Node: 1,
					Core: DataConfig{
						Start: 0,
						Bytes: 0,
					},
				},
			},
			TCP: ModbusTCPConfig{
				Address: address,
				Port:    502,
				Read:    60,
				Write:   15,
				Control: 15,
				Core: DataConfig{
					Start: 0,
					Bytes: 0,
				},
			},
		},
		UA: UAConfig{
			Port:      54000,
			Namespace: 1,
		},
		DB: DBConfig{
			Address:  address,
			Proto:    "tcp",
			Port:     "5432",
			Timeout:  10000,
			Base:     os.Getenv("IRO_POSTGRES_DATABASE"),
			User:     os.Getenv("IRO_POSTGRES_USER"),
			Password: os.Getenv("IRO_POSTGRES_PASSWORD"),
		},
		HTTP: HTTPConfig{
			Address:  address,
			Port:     "80",
			Read:     60,
			Write:    15,
			Shutdown: 120,
			Root:     "./",
		},
	}
}

// Load - загрузка данных конфигурации
func (o *SystemConfig) Load(path *string) error {
	file, err := os.ReadFile(*path)
	if err == nil {
		err = yaml.Unmarshal(file, o)
	}
	return err
}

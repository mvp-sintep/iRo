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
	Pause    int    `yaml:"pause"`
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
	Core    DataConfig `yaml:"core"`
}

// ModbusConfig - каталог записей конфигурации MODBUS
type ModbusConfig struct {
	RTU []ModbusRTUConfig `yaml:"rtu"`
	TCP ModbusTCPConfig   `yaml:"tcp"`
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

// SystemConfiguration - все данные настройки
type SystemConfiguration struct {
	COM    []COMPortConfig `yaml:"com"`
	Modbus ModbusConfig    `yaml:"modbus"`
	HTTP   HTTPConfig      `yaml:"http"`
}

// New - создание записи конфигурации
func New() *SystemConfiguration {

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
	return &SystemConfiguration{
		COM: []COMPortConfig{
			{
				File:     "/dev/ttyS0",
				BaudRate: 115200,
				DataBits: 8,
				Parity:   "N",
				StopBits: 1,
				Timeout:  1,
				Pause:    5,
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
				Core: DataConfig{
					Start: 0,
					Bytes: 0,
				},
			},
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
func (o *SystemConfiguration) Load(path *string) error {
	file, err := os.ReadFile(*path)
	if err == nil {
		err = yaml.Unmarshal(file, o)
	}
	return err
}

package config

type Config struct {
	ServiceName string `json:"service_name"`
	MySQL       MySQL  `json:"mysql"`
	JWT         JWT    `json:"jwt"`
}

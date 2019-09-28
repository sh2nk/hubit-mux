package config

//Config - храним тут параметры конфигурации
type Config struct {
	Addr    string `json:"addr"`
	Device  string `json:"device"`
	WB      bool   `json:"wb"`
	ServURL string `json:"serv-url"`
	Format  string `json:"image-format"`
	Width   uint32 `json:"width"`
	Height  uint32 `json:"height"`
}

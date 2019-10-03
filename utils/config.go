package utils

//Configuration - тип параметров конфигурации
type Configuration struct {
	Addr      string `json:"addr"`
	Device    string `json:"device"`
	WB        bool   `json:"wb"`
	StreamURL string `json:"stream-url"`
	Format    string `json:"image-format"`
	Width     uint32 `json:"width"`
	Height    uint32 `json:"height"`
	Resize    uint32 `json:"resize-level"`
}

//Config - храним тут параметры конфигурации
var Config Configuration

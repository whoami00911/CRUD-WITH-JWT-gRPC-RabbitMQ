package domain

import (
	"sync"
	"time"
)

type AssetsResponse struct {
	Timestamp int64                    `json:"timestamp"`
	Cache     map[*AssetData]time.Time `json:"data"`
	Mu        sync.Mutex
}
type AssetResponse struct {
	Asset AssetData `json:"data"`
	Mu    sync.Mutex
}
type AssetData struct {
	IPAddress            string `json:"ipAddress"`            // IP-адрес актива
	IsPublic             bool   `json:"isPublic"`             // Является ли IP-адрес публичным
	IPVersion            int    `json:"ipVersion"`            // Версия IP (4 или 6)
	IsWhitelisted        bool   `json:"isWhitelisted"`        // Находится ли IP-адрес в белом списке
	AbuseConfidenceScore int    `json:"abuseConfidenceScore"` // Оценка злоупотребления (от 0 до 100)
	CountryCode          string `json:"countryCode"`          // Код страны
	CountryName          string `json:"countryName"`          // Название страны
	UsageType            string `json:"usageType"`            // Тип использования (например, дата-центр/веб-хостинг)
	ISP                  string `json:"isp"`                  // Интернет-провайдер (ISP)
	IsTor                bool   `json:"isTor"`                // Является ли IP узлом выхода Tor
	IsCache              bool
	IsDb                 bool `json:"isFromDB"`
}

func NewAsset() *AssetResponse {
	return &AssetResponse{}
}

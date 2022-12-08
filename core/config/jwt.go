package config

type JWT struct {
	SigningKey  string `json:"signing_key"`  // jwt签名
	ExpiresTime string `json:"expires_time"` // 过期时间
	BufferTime  string `json:"buffer_time"`  // 缓冲时间
	Issuer      string `json:"issuer"`       // 签发者
}

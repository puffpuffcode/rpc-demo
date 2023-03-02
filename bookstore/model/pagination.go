// 基于 cursor 的分页

package model

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Token string

type Page struct {
	NextID string `json:"next_id"`
	NextTimeAtUTC int64 `json:"next_time_at_utc"`
	PageSize int64 `json:"page_size"`
}

// Token 编码
func (p *Page) Encode() Token {
	b, err := json.Marshal(p)
	if err != nil {
		return Token("")
	}
	return Token(base64.StdEncoding.EncodeToString(b))
}

// 解析 Token
func (token Token) Decode() Page {
	var page Page
	if len(token) == 0 {
		return page
	}
	b, err := base64.StdEncoding.DecodeString(string(token))
	if err != nil {
		return page
	}
	if err = json.Unmarshal(b,&page); err != nil {
		return page
	}
	return page
}

// 判断 Token 是否 有效
func (p *Page) IsInVaild() bool {
	return p.NextID == "" || p.NextTimeAtUTC == 0 || p.NextTimeAtUTC > time.Now().Unix() || p.PageSize == 0
}
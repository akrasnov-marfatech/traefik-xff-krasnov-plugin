package traefik_xff_krasnov_plugin

import (
	"context"
	"net/http"
	"strings"
)

type Config struct {
	// Если true — при отсутствии XFF ставим пустую строку (""), иначе удаляем заголовок
	KeepEmpty bool `json:"keepEmpty,omitempty"`
	// Если true — оставляем только первый IP до запятой, иначе весь заголовок как прислал клиент
	OnlyFirst bool `json:"onlyFirst,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		KeepEmpty: true,
		OnlyFirst: true,
	}
}

type Middleware struct {
	next http.Handler
	cfg  *Config
}

func New(_ context.Context, next http.Handler, cfg *Config, _ string) (http.Handler, error) {
	return &Middleware{next: next, cfg: cfg}, nil
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	xff := req.Header.Get("X-Forwarded-For")

	if xff == "" {
		if m.cfg.KeepEmpty {
			req.Header.Set("X-Forwarded-For", "")
		} else {
			req.Header.Del("X-Forwarded-For")
		}
	} else {
		if m.cfg.OnlyFirst {
			first := strings.TrimSpace(strings.Split(xff, ",")[0])
			req.Header.Set("X-Forwarded-For", first)
		} else {
			req.Header.Set("X-Forwarded-For", xff)
		}
	}

	// при желании можно дублировать «чистое» значение:
	// req.Header.Set("X-Client-IP", req.Header.Get("X-Forwarded-For"))

	m.next.ServeHTTP(rw, req)
}

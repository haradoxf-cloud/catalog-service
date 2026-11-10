package rprocessor

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/haradoxf-cloud/catalog-service/internal/app/config/section"
	rhandler "github.com/haradoxf-cloud/catalog-service/internal/app/handler/http"
)

type httpProc struct {
	server http.Server
	addr   string
}

func NewHTTP(hHealth rhandler.Health, cfg section.ProcessorWebServer) *httpProc {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(handlerNotFound)

	vGenericRegHealthCheck(r, hHealth)
	_ = r.Walk(func(rt *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		p, _ := rt.GetPathTemplate()
		m, _ := rt.GetMethods()
		if p != "" && len(m) > 0 {
			log.Printf("[ROUTER] %s -> %s", m, p)
		}
		return nil
	})

	p := httpProc{addr: fmt.Sprintf(":%d", cfg.ListenPort)}
	p.server.Addr = p.addr
	p.server.Handler = r

	return &p
}

func (p *httpProc) Serve() error {
	log.Printf("Starting HTTP server on %s", p.addr)
	return p.server.ListenAndServe()
}

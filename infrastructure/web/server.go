package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web/handler"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web/router"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init(h handler.AppHandler) {
	s.router = router.NewMuxRouter(h)
}

func (s *Server) Run(port int) {
	fmt.Printf("Server running at http://loacalhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), s.router)
	return
}

package stat

import (
	"URLshortener/configs"
	"URLshortener/pkg/middleware"
	"URLshortener/pkg/res"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	Config         *configs.Config
	StatRepository *StatRepository
}

type StatHandler struct {
	Config         *configs.Config
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		Config:         deps.Config,
		StatRepository: deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}
		stats := handler.StatRepository.GetStat(by, from, to)
		res.Json(w, stats, http.StatusOK)
	}
}

package save

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
	"net/http"
	"short_url/internal/lib/logger/sl"
	"strconv"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

type URLSaver interface {
	SaveUrl(url, alias string) (int64, error) 
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),				
		)
		
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			
			render.JSON(w, r, Response{
				Status: "error",
				Error: "failed to decode request",
			})
			
			return
		}
		
		log.Info("request decoded", slog.Any("request", req))
//		log.Info("saving url", slog.String("url", req.URL), slog.String("alias", req.Alias))
		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))
			
			render.JSON(w, r, Response{
				Status: "error",
				Error: "invalid request",
			})
			
			return
		}
		
		alias := req.Alias
		if alias == "" {
			alias = "Test alias" + middleware.GetReqID(r.Context())
		}
		
		id, err := urlSaver.SaveUrl(req.URL, alias)
		if err != nil {
			
			fmt.Print("SMTH HAPPENED")
			log.Info("Smth happened")
			
			render.JSON(w, r, Response{
				Error: "smth happened",
				Status: "error",
			})
			
			return
		}
		
		render.JSON(w, r, Response{
			Status: "OKEY",
			Error: strconv.FormatInt(id, 10),
		})
	}
}


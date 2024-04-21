package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/DanglingDynamo/chronotube/internal/models"
	"github.com/DanglingDynamo/chronotube/internal/services"
	"github.com/DanglingDynamo/chronotube/internal/utils"
)

type SearchHandler interface {
	FetchVideosPaginated(writer http.ResponseWriter, req *http.Request)
}

type searchHandler struct {
	service services.SearchService
}

func NewSearchHandler(service services.SearchService) SearchHandler {
	return &searchHandler{
		service: service,
	}
}

func (sh *searchHandler) FetchVideosPaginated(writer http.ResponseWriter, req *http.Request) {
	var request models.PaginatedVideoRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		utils.WriteJSON(writer, http.StatusBadRequest, map[string]string{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	videos, nextPage, err := sh.service.QueryVideos(req.Context(), request)
	if err != nil {
		utils.WriteJSON(writer, http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"videos":    videos,
			"next_page": nextPage,
		},
	})
}

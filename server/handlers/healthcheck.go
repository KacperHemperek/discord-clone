package handlers

import (
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

func HandleHealthcheck() utils.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"status": "ok"})
	}
}

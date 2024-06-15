package handlers

import (
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

func HandleHealthcheck() utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request, _ *utils.APIContext) error {
		return utils.WriteJson(w, http.StatusOK, utils.JSON{"status": "ok"})
	}
}

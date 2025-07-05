package controllers

import (
	services "lotto-backend-api/services"
	utils "lotto-backend-api/utils"

	"net/http"

	"gorm.io/gorm"
)

// GetUsers handles GET /users
func RegisterController(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	response := services.RegisterUser(DB, r)
	utils.RespondWithJSON(w, http.StatusOK, response)
}

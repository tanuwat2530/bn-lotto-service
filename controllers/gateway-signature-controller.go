package controllers

import (
	services "lotto-backend-api/services"
	utils "lotto-backend-api/utils"

	"net/http"

	"gorm.io/gorm"
)

// GetUsers handles GET /users
func GatewaySignatureController(DB *gorm.DB, w http.ResponseWriter, r *http.Request) {
	response := services.GatewaySignature()
	utils.RespondWithJSON(w, http.StatusOK, response)
}

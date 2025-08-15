package routes

import (
	handlers "lotto-backend-api/handlers"
	"net/http"

	"gorm.io/gorm"
)

// Import the new handlers package

// SetupRoutes now accepts the database connection as an argument.
// It creates the handlers, injects the db, and sets up the routes.
func SetupRoutes(db *gorm.DB) http.Handler {
	// Use a powerful router like chi or gin in a real app.
	// For this example, http.NewServeMux is clear and simple.
	mux := http.NewServeMux()

	// --- Dependency Injection into Handlers ---
	// Create an instance of our application handlers, passing the db connection.
	lottoHandlers := &handlers.LottoHandler{DB: db}

	// --- Register Routes ---
	// Now, we map URL paths to the methods of our handler instance.
	// Each of these handlers now has access to the database via `lottoHandlers.DB`.
	mux.HandleFunc("/lotto-api/register", lottoHandlers.Register)
	mux.HandleFunc("/lotto-api/login", lottoHandlers.Login)
	mux.HandleFunc("/lotto-api/credit", lottoHandlers.CreditBalance)
	mux.HandleFunc("/lotto-api/bet", lottoHandlers.Bet)
	mux.HandleFunc("/lotto-api/history", lottoHandlers.History)

	//mux.HandleFunc("/gateway-api/payment-channel", lottoHandlers.PaymentChannel)
	//mux.HandleFunc("/gateway-api/bank-provider", lottoHandlers.BankProvider)
	//mux.HandleFunc("/gateway-api/payment-status", lottoHandlers.PaymentNoti)
	mux.HandleFunc("/gateway-api/pay-in", lottoHandlers.PayIn)
	mux.HandleFunc("/gateway-api/pay-out", lottoHandlers.PayOut)
	mux.HandleFunc("/gateway-api/order-noti", lottoHandlers.PaymentNoti)

	mux.HandleFunc("/lotto-api/", HomeHandler)
	return mux
}

// HomeHandler for root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the LOTTO backend API power by GoLang ^_^"))
}

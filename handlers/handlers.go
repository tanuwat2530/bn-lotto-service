package handlers

import (
	controllers "lotto-backend-api/controllers"
	"net/http"

	"gorm.io/gorm"
)

// LottoHandler holds dependencies for our request handlers.
// In this case, it's the GORM database connection pool.
type LottoHandler struct {
	DB *gorm.DB
}

// GetLottoResults is a method on LottoHandler. Because it's a method,
// it has access to all fields of the LottoHandler struct, including `h.DB`.
func (h *LottoHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.RegisterController(h.DB, w, r)
}

func (h *LottoHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.LoginController(h.DB, w, r)
}

func (h *LottoHandler) Bet(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.BetController(h.DB, w, r)
}

func (h *LottoHandler) History(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.HistoryController(h.DB, w, r)
}

func (h *LottoHandler) CreditBalance(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.CreditBalanceController(h.DB, w, r)
}

func (h *LottoHandler) PaymentChannel(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//controllers.PaymentChannelController(h.DB, w, r)
}

func (h *LottoHandler) BankProvider(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//controllers.BankProviderController(h.DB, w, r)
}

func (h *LottoHandler) PayIn(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.PayInController(h.DB, w, r)
}

func (h *LottoHandler) PayOut(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.PayOutController(h.DB, w, r)
}

func (h *LottoHandler) PaymentNoti(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.OrderNotiController(h.DB, w, r)
}

func (h *LottoHandler) PromtpayCreditNoti(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.PromtpayNotiController(h.DB, w, r)
}

func (h *LottoHandler) Promtpay(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controllers.QrPromtpayController(h.DB, w, r)
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/domain"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type PaymentService struct {
	paymentRepo  domain.PaymentRepository
	customerRepo domain.CustomerRepository
	authRepo     domain.AuthRepository
	trm          trm.Manager
	cfg          config.Choco
}

func NewPaymentService(
	paymentRepository domain.PaymentRepository,
	customerRepository domain.CustomerRepository,
	authRepository domain.AuthRepository,
	trm trm.Manager,
	cfg config.Choco,
) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepository,
		customerRepo: customerRepository,
		authRepo:     authRepository,
		trm:          trm,
		cfg:          cfg,
	}
}

type TransactionListResponseData struct {
	ErrorCode int    `json:"error_code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Data      struct {
		Pagination struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			TotalItems int `json:"total_items"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
		Items []struct {
			UserID int64 `json:"user_id"`
			// Other fields omitted for brevity
		} `json:"items"`
	} `json:"data"`
}

func (s *PaymentService) fetchUniqueUserIDs(ctx context.Context, baseURL string) ([]int64, error) {
	page := 1
	userIDMap := make(map[int64]struct{}) // Use map to store unique user IDs
	client := &http.Client{}

	for {
		// Construct the paginated URL
		url := fmt.Sprintf("%s&page=%d", baseURL, page)

		token := s.cfg.ChocoToken
		//token, err := getAccessToken(ctx, s.authRepo, s.cfg)
		//if err != nil {
		//	return nil, fmt.Errorf("failed to get access token: %w", err)
		//}

		body, err := sendChocoRequest(ctx, client, url, token)
		if err != nil && err.Error() == "unexpected status code: 401" {
			token, err = fetchNewAccessToken(ctx, s.authRepo, client, s.cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch new access token: %w", err)
			}
			body, err = sendChocoRequest(ctx, client, url, token)
			if err != nil {
				return nil, fmt.Errorf("failed to send choco request: %w", err)
			}
		}

		var responseData TransactionListResponseData
		if err := json.Unmarshal(body, &responseData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		// Collect unique user IDs from this page
		for _, item := range responseData.Data.Items {
			userIDMap[item.UserID] = struct{}{}
		}

		// Check if there are more pages
		if page >= responseData.Data.Pagination.TotalPages {
			break
		}

		page++
		//time.Sleep(2 * time.Second)
	}

	// Convert the map keys to a slice
	userIDs := make([]int64, 0, len(userIDMap))
	for userID := range userIDMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (s *PaymentService) FetchPayments(ctx context.Context, terminals string) error {
	customerService := NewCustomerService(s.customerRepo, s.authRepo, s.trm, s.cfg)

	var endDate = time.Now().Format("2006-01-02") + "+23:59:59"
	var startDate = time.Now().AddDate(0, 0, -2).Format("2006-01-02") + "+00:00:00"

	baseURL := "https://api-proxy.choco.kz/acl/proxy?proxy_path=reports/merchant/transactions&filials=" +
		terminals +
		"&types=pay,refund&start_date=" + startDate + "&end_date=" + endDate

	fmt.Println("fetching user IDs")
	userIDs, err := s.fetchUniqueUserIDs(ctx, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch unique user IDs: %w", err)
	}

	fmt.Println(string(rune(len(userIDs))) + " customers found")

	for _, userID := range userIDs {
		fmt.Println("fetching customer info for user ID: " + strconv.FormatInt(userID, 10))
		_, err := customerService.FetchCustomerInfo(ctx, domain.CustomerID(userID), terminals)
		if err != nil {
			return fmt.Errorf("failed to fetch customer info: %w", err)
		}

		fmt.Println("fetching user payments for user ID: " + strconv.FormatInt(userID, 10))
		if err := s.fetchUserPayments(ctx, terminals, userID, startDate, endDate); err != nil {
			return fmt.Errorf("failed to fetch user payments: %w", err)
		}

		// sleep for 3 sec
		//time.Sleep(3 * time.Second)
	}

	return nil
}

type PaymentHistoryResponseData struct {
	JSONAPI struct {
		Version string `json:"version"`
	} `json:"jsonapi"`
	Meta struct {
		Page struct {
			CurrentPage int `json:"currentPage"`
			LastPage    int `json:"lastPage"`
			PerPage     int `json:"perPage"`
			Total       int `json:"total"`
		} `json:"page"`
	} `json:"meta"`
	Data []PHDataItem `json:"data"`
}

type PHDataItem struct {
	ID         int64         `json:"id"`
	Type       string        `json:"type"`
	Attributes []PHAttribute `json:"attributes"`
}

type PHAttribute struct {
	Transaction struct {
		ID             int64   `json:"id"`
		CreatedBy      int64   `json:"created_by"`
		Type           string  `json:"type"`
		Amount         float64 `json:"amount"`
		DiscountAmount float64 `json:"discount_amount"`
		CreatedAt      string  `json:"created_at"`
	} `json:"transaction"`
	Location struct {
		Title     string `json:"title"`
		PartnerID string `json:"partner_id"`
	} `json:"location"`
	Review interface{} `json:"review"`
}

func (s *PaymentService) fetchUserPayments(
	ctx context.Context,
	terminals string,
	userId int64,
	startDate string,
	endDate string,
) error {
	client := &http.Client{}
	page := 1

	baseUrl := "https://api-proxy.choco.kz/analytics/v1/customer/" +
		strconv.FormatInt(userId, 10) +
		"/payment-history?terminals=" +
		terminals +
		"&filter[start_date]=" + startDate + "&filter[end_date]=" + endDate

	for {
		// Construct the paginated URL
		url := fmt.Sprintf("%s&page=%d", baseUrl, page)

		token := s.cfg.ChocoToken
		//token, err := getAccessToken(ctx, s.authRepo, s.cfg)
		//if err != nil {
		//	return fmt.Errorf("failed to get access token: %w", err)
		//}

		body, err := sendChocoRequest(ctx, client, url, token)
		if err != nil {
			if err.Error() == "unexpected status code: 401" {
				token, err = fetchNewAccessToken(ctx, s.authRepo, client, s.cfg)
				if err != nil {
					return fmt.Errorf("failed to fetch new access token: %w", err)
				}
				body, err = sendChocoRequest(ctx, client, url, token)
				if err != nil {
					return fmt.Errorf("failed to send choco request: %w", err)
				}
			} else {
				return fmt.Errorf("failed to send choco request: %w", err)
			}

		}

		// Process the response body
		var responseData PaymentHistoryResponseData
		if err := json.Unmarshal(body, &responseData); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		for _, item := range responseData.Data {
			err := s.storePayment(ctx, item)
			if err != nil {
				return fmt.Errorf("failed to store payment: %w", err)
			}
		}

		// Check if there are more pages
		if page >= responseData.Meta.Page.LastPage {
			break
		}
		page++
	}

	return nil
}

func (s *PaymentService) storePayment(ctx context.Context, item PHDataItem) error {
	payment := domain.Payment{
		ID:                domain.PaymentID(item.Attributes[0].Transaction.ID),
		CreatedBy:         item.Attributes[0].Transaction.CreatedBy,
		Type:              item.Attributes[0].Transaction.Type,
		Amount:            int64(item.Attributes[0].Transaction.Amount),
		DiscountAmount:    int64(item.Attributes[0].Transaction.DiscountAmount),
		CreatedAt:         item.Attributes[0].Transaction.CreatedAt,
		LocationTitle:     item.Attributes[0].Location.Title,
		LocationPartnerID: item.Attributes[0].Location.PartnerID,
	}

	exists, err := s.paymentRepo.ExistsById(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("failed to check if payment exists: %w", err)
	}

	if exists {
		return nil
	}
	_, err = s.paymentRepo.Create(ctx, &payment)
	if err != nil {
		return fmt.Errorf("failed to store payment: %w", err)
	}

	return nil
}

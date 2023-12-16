// Package docs Gophermart API
//
// Documentation for Gophermart API
//
//		Schemes: http
//		Host: localhost:8080
//		BasePath: /api/user
//		Version: 1.0.0
//
//		Consumes:
//		- application/json
//		- text/plain
//
//		Produces:
//		- application/json
//
//	 Security:
//	 - api_key
//
//	 SecurityDefinitions:
//	 api_key:
//	   type: apiKey
//	   name: Authorization
//	   in: header
//
// swagger:meta
package docs

import (
	"github.com/leonf08/gophermart.git/internal/models"
)

// swagger:route POST /register auth userSignUp
// Register a new user.
// consumes:
// - application/json
// responses:
//   200: noContentResponse
//   400: errorResponse
//   409: errorResponse
//   500: errorResponse

// swagger:route POST /login auth userLogIn
// Log in a user.
// consumes:
// - application/json
// responses:
//   200: noContentResponse
//   400: errorResponse
//   401: errorResponse
//   500: errorResponse

// swagger:route POST /orders orders uploadOrder
// Upload an order.
// consumes:
// - text/plain
// security:
//   api_key:
// responses:
//   200: noContentResponse
//   202: noContentResponse
//   400: errorResponse
//   401: errorResponse
//   409: errorResponse
//   422: errorResponse
//   500: errorResponse

// swagger:route GET /orders orders getOrders
// Get orders.
// security:
//   api_key:
// responses:
//   200: getOrdersResponse
//   204: noContentResponse
//   401: errorResponse
//   500: errorResponse

// swagger:route GET /balance balance getUserBalance
// Get user balance.
// security:
//   api_key:
// responses:
//   200: getBalanceResponse
//   401: errorResponse
//   500: errorResponse

// swagger:route POST /balance/withdraw balance withdraw
// Withdraw money from the user balance.
// consumes:
// - application/json
// security:
//   api_key:
// responses:
//   200: noContentResponse
//   400: errorResponse
//   401: errorResponse
//   402: errorResponse
//   422: errorResponse
//   500: errorResponse

// swagger:route GET /withdrawals balance getWithdrawals
// Get withdrawals.
// security:
//   api_key:
// responses:
//   200: getOrdersResponse
//   204: noContentResponse
//   401: errorResponse
//   500: errorResponse

// swagger:parameters userSignUp userLogIn
type userSignUpRequest struct {
	// in: body
	Body *models.User
}

// swagger:parameters uploadOrder
type uploadOrderRequest struct {
	// in: body
	Body string
}

// swagger:parameters withdraw
type withdrawRequest struct {
	// in: body
	Body struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}
}

// noContentResponse is a response body when content is empty.
// swagger:response noContentResponse
type noContentResponse struct{}

// getOrdersResponse is a response body for the getOrders handler when the input is valid.
// swagger:response getOrdersResponse
type getOrdersResponse struct {
	// in: body
	Body []models.Order
}

// getBalanceResponse is a response body for the getUserBalance handler when the input is valid.
// swagger:response getBalanceResponse
type getBalanceResponse struct {
	// in: body
	Body *models.UserAccount
}

// getWithdrawalsResponse is a response body for the getWithdrawals handler when the input is valid.
// swagger:response getWithdrawalsResponse
type getWithdrawalsResponse struct {
	// in: body
	Body []models.Withdrawal
}

// errorResponse is a response body for the userSignUp handler when the input is invalid.
// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Err string
}

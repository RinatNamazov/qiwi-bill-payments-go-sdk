/*****************************************************************************
 *
 *  PROJECT:        qiwi-bill-payments-go-sdk
 *  LICENSE:        See LICENSE in the top level directory
 *  FILE:           types.go
 *  DESCRIPTION:    Types
 *  COPYRIGHT:      (c) 2021 RINWARES <rinwares.com>
 *  AUTHOR:         Rinat Namazov <rinat.namazov@rinwares.com>
 *
 *****************************************************************************/

package qiwi

import (
	"fmt"
	"strconv"
	"time"
)

// ErrorResponse is the API error response.
type ErrorResponse struct {
	ServiceName string    `json:"serviceName"` // The service name.
	ErrorCode   string    `json:"errorCode"`   // The error code.
	Description string    `json:"description"` // The description.
	UserMessage string    `json:"userMessage"` // The user message.
	DateTime    time.Time `json:"dateTime"`    // The dateTime.
	TraceId     string    `json:"traceId"`     // The trace ID.
}

// Error returns a formatted error string.
func (m ErrorResponse) Error() string {
	return "ServiceName: " + m.ServiceName + ", ErrorCode: " + m.ErrorCode + ", Description: " + m.Description +
		", UserMessage: " + m.UserMessage + ", DateTime: " + m.DateTime.Format(dateTimeFormat) + ", TraceId: " + m.TraceId
}

// MoneyAmount is the invoice amount info.
type MoneyAmount struct {
	Value    string `json:"value"`    // The invoice amount value.
	Currency string `json:"currency"` // The invoice currency value.
}

// SetValueString sets the amount from a string.
func (m *MoneyAmount) SetValueString(value string) error {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.Value = fmt.Sprintf("%.2f", v)
	return nil
}

// SetValueString sets the amount from a float.
func (m *MoneyAmount) SetValueNumber(value float64) {
	m.Value = fmt.Sprintf("%.2f", value)
}

// GetValueNumber returns the amount in float.
func (m *MoneyAmount) GetValueNumber() float64 {
	v, _ := strconv.ParseFloat(m.Value, 64)
	return v
}

// Customer is the customer's info.
type Customer struct {
	Email   string `json:"email"`   // The client's e-mail.
	Account string `json:"account"` // The client's identifier in merchant's system.
	Phone   string `json:"phone"`   // The phone number to which invoice issued.
}

// ResponseStatus is the invoice status info.
type ResponseStatus struct {
	Value           string    `json:"value"`           // The status value.
	ChangedDateTime time.Time `json:"changedDateTime"` // The status refresh dateTime.
}

// CustomFields is the invoice additional data.
type CustomFields struct {
	ApiClient        string `json:"apiClient"`        // // The API client name.
	ApiClientVersion string `json:"apiClientVersion"` // The API client version.
	ThemeCode        string `json:"themeCode"`        // The style theme code.
}

// BillResponse is the invoice response.
type BillResponse struct {
	SiteId             string         `json:"siteId"`             // The merchant’s site identifier in API.
	BillId             string         `json:"billId"`             // The unique invoice identifier in the merchant’s system.
	Amount             MoneyAmount    `json:"amount"`             // The invoice amount info.
	Status             ResponseStatus `json:"status"`             // The invoice status info.
	Comment            string         `json:"comment"`            // The comment to the invoice.
	Customer           Customer       `json:"customer"`           // The customer info.
	CreationDateTime   time.Time      `json:"creationDateTime"`   // The dateTime of the invoice creation.
	ExpirationDateTime time.Time      `json:"expirationDateTime"` // The expiration date of the pay form link.
	PayUrl             string         `json:"payUrl"`             // The pay form link.
	CustomFields       CustomFields   `json:"customFields"`       // The invoice additional data.
}

// Notification is the invoice payment notification.
type Notification struct {
	BillResponse
	Version string `json:"version"` // The notification version.
}

// CreateBillInfo is create issue info.
type CreateBillInfo struct {
	BillId             string      `json:"billId"`             // The unique invoice identifier in merchant's system.
	Amount             MoneyAmount `json:"amount"`             // The invoice amount info.
	Comment            string      `json:"comment"`            // The invoice commentary.
	ExpirationDateTime time.Time   `json:"expirationDateTime"` // The invoice due date.
	Customer           Customer    `json:"customer"`           // The customer's info.
	SuccessUrl         string      `json:"successUrl"`         // The URL to which the client will be redirected in case of successful payment.
	ThemeCode          string      `json:"-"`                  // The style theme code.
}

// CreateBillRequest is create issue request.
type CreateBillRequest struct {
	Amount             MoneyAmount  `json:"amount"`             // The invoice amount witch currency.
	Comment            string       `json:"comment"`            // The invoice commentary.
	ExpirationDateTime time.Time    `json:"expirationDateTime"` // The invoice expiration date.
	Customer           Customer     `json:"customer"`           // The customer's info.
	CustomFields       CustomFields `json:"customFields"`       // The invoice additional data.
}

// GetCreateBillRequest returns a create invoice request.
func (m *CreateBillInfo) GetCreateBillRequest() CreateBillRequest {
	return CreateBillRequest{
		Amount:             m.Amount,
		Comment:            m.Comment,
		ExpirationDateTime: m.ExpirationDateTime,
		Customer:           m.Customer,
		CustomFields: CustomFields{
			ThemeCode: m.ThemeCode,
		},
	}
}

// PaymentInfo is the invoice data are put in Pay Form URL.
type PaymentInfo struct {
	PublicKey  string      `json:"publicKey"`  // The merchant public key.
	Amount     MoneyAmount `json:"amount"`     // The invoice amount.
	BillId     string      `json:"billId"`     // Unique invoice identifier in merchant’s system.
	SuccessUrl string      `json:"successUrl"` // The URL to which the client will be redirected in case of successful payment.
	ThemeCode  string      `json:"-"`          // The style theme code..
}

// RefundBillRequest is the refund request.
type RefundBillRequest struct {
	Amount MoneyAmount `json:"amount"` // The refund amount.
}

// RefundResponse is the refund response.
type RefundResponse struct {
	Amount   MoneyAmount `json:"amount"`   // The invoice amount.
	DateTime time.Time   `json:"dateTime"` // The dateTime of refund processing.
	RefundId string      `json:"refundId"` // Unique refund identifier in merchant’s system.
	Status   string      `json:"status"`   // The refund status.
}

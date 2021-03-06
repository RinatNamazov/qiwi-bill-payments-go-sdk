/*****************************************************************************
 *
 *  PROJECT:        qiwi-bill-payments-go-sdk
 *  LICENSE:        See LICENSE in the top level directory
 *  FILE:           qiwi.go
 *  DESCRIPTION:    Main file
 *  COPYRIGHT:      (c) 2021 RINWARES <rinwares.com>
 *  AUTHOR:         Rinat Namazov <rinat.namazov@rinwares.com>
 *
 *****************************************************************************/

package qiwi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	clientName    = "go_sdk"
	clientVersion = "1.1.0"
)

// The API base URL.
const (
	billsURL   = "https://api.qiwi.com/partner/bill/v1/bills/"
	paymentURL = "https://oplata.qiwi.com/create"
)

// The API dateTime format.
const dateTimeFormat = "2006-01-02T15:04:05.999999-07:00"

// Invoice Payment Statuses
const (
	WAITING  = "WAITING"  // Invoice issued awaiting for payment.
	PAID     = "PAID"     // Invoice paid.
	REJECTED = "REJECTED" // Invoice rejected by customer.
	EXPIRED  = "EXPIRED"  // Invoice expired. Invoice not paid.
)

// The refund status enum.
const (
	PARTIAL = "PARTIAL" // The partial refund of the invoice amount.
	FULL    = "FULL"    // The full refund of the invoice amount.
)

// HTTPClient is HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// QiwiBillPaymentsAPI for rest v3.
type QiwiBillPaymentsAPI struct {
	httpClient HTTPClient
	secretKey  string
}

// NewQiwiBillPaymentsAPI creates a new QiwiBillPaymentsAPI instance.
func NewQiwiBillPaymentsAPI(secretKey string) *QiwiBillPaymentsAPI {
	return NewQiwiBillPaymentsAPIWithClient(secretKey, &http.Client{})
}

// NewQiwiBillPaymentsAPIWithClient creates a new QiwiBillPaymentsAPI instance
// and allows you to pass a http.Client.
func NewQiwiBillPaymentsAPIWithClient(secretKey string, httpClient HTTPClient) *QiwiBillPaymentsAPI {
	return &QiwiBillPaymentsAPI{
		httpClient: httpClient,
		secretKey:  "Bearer " + secretKey,
	}
}

// makeRequest makes a request.
func (m *QiwiBillPaymentsAPI) makeRequest(url, method string, body interface{}, resp interface{}) error {
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, billsURL+url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", m.secretKey)

	r, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	respData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if r.StatusCode != 200 {
		var errRes ErrorResponse
		err = json.Unmarshal(respData, &errRes)
		if err != nil {
			return err
		}
		return errRes
	} else if resp != nil {
		err = json.Unmarshal(respData, &resp)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetSecretKey sets new secret key.
func (m *QiwiBillPaymentsAPI) SetSecretKey(secretKey string) {
	m.secretKey = "Bearer " + secretKey
}

// CreatePaymentForm creating checkout link.
func (m *QiwiBillPaymentsAPI) CreatePaymentForm(paymentInfo PaymentInfo) string {
	var v url.Values
	v.Add("amount", paymentInfo.Amount.Value)
	v.Add("publicKey", paymentInfo.PublicKey)
	v.Add("billId", paymentInfo.BillId)
	v.Add("successUrl", paymentInfo.SuccessUrl)
	v.Add("customFields[apiClient]", clientName)
	v.Add("customFields[apiClientVersion]", clientVersion)
	if paymentInfo.ThemeCode != "" {
		v.Add("customFields[themeCode]", paymentInfo.ThemeCode)
	}
	return paymentURL + "?" + v.Encode()
}

// GetBillInfo getting bill info.
func (m *QiwiBillPaymentsAPI) GetBillInfo(billId string) (*BillResponse, error) {
	var resp BillResponse
	err := m.makeRequest(billId, "GET", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateBill creating bill.
func (m *QiwiBillPaymentsAPI) CreateBill(billInfo CreateBillInfo) (*BillResponse, error) {
	billRequest := billInfo.GetCreateBillRequest()
	billRequest.CustomFields.ApiClient = clientName
	billRequest.CustomFields.ApiClientVersion = clientVersion

	var resp BillResponse
	err := m.makeRequest(billInfo.BillId, "PUT", billRequest, &resp)
	if err != nil {
		return nil, err
	}

	if billInfo.SuccessUrl != "" {
		resp.PayUrl = resp.PayUrl + "&successUrl=" + url.QueryEscape(billInfo.SuccessUrl)
	}

	return &resp, nil
}

// CancelBill cancelling unpaid bill.
func (m *QiwiBillPaymentsAPI) CancelBill(billId string) (*BillResponse, error) {
	var resp BillResponse
	err := m.makeRequest(billId+"/reject", "POST", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRefundInfo getting refund info.
// Method is not available for individuals.
func (m *QiwiBillPaymentsAPI) GetRefundInfo(billId, refundId string) (*RefundResponse, error) {
	var resp RefundResponse
	err := m.makeRequest(billId+"/refunds/"+refundId, "GET", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Refund paid bill.
// Method is not available for individuals.
func (m *QiwiBillPaymentsAPI) Refund(billId, refundId string, amount MoneyAmount) (*RefundResponse, error) {
	req := RefundBillRequest{Amount: amount}

	var resp RefundResponse
	err := m.makeRequest(billId+"/refunds/"+refundId, "PUT", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

# qiwi-bill-payments-go-sdk

Golang SDK for QIWI Payments.

## Installation

```bash
go get github.com/RinatNamazov/qiwi-bill-payments-go-sdk
```

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/RinatNamazov/qiwi-bill-payments-go-sdk"
	"github.com/google/uuid"
)

func main() {
	qiwiApi := qiwi.NewQiwiBillPaymentsAPI("SecretKey")

	billInfo := qiwi.CreateBillInfo{
		BillId:             uuid.New().String(),
		Comment:            "Donation for coffee.",
		ExpirationDateTime: time.Now().Add(time.Hour),
	}
	billInfo.Amount.SetValueNumber(100)
	billInfo.Amount.Currency = "RUB"

	billResp, err := qiwiApi.CreateBill(billInfo)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Amount: " + billInfo.Amount.Value + " " + billInfo.Amount.Currency +
			"\nBill ID: " + billInfo.BillId +
			"\nFor payment follow the link:\n" + billResp.PayUrl)
	}
}
```

## License
The code in this repository is licensed under MIT.
```
Copyright (c) 2021 RINWARES, Rinat Namazov

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

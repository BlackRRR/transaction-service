package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type TestRequest struct {
	ClientID int64 `json:"client_id"`
	Amount   int64 `json:"amount"`
}

func Test_transaction(t *testing.T) {
	var AllReqs []TestRequest

	req1 := TestRequest{
		ClientID: 1215125125,
		Amount:   300,
	}

	req2 := TestRequest{
		ClientID: 2412412412,
		Amount:   6636,
	}

	req3 := TestRequest{
		ClientID: 4444444444,
		Amount:   555,
	}

	req4 := TestRequest{
		ClientID: 1215125125,
		Amount:   15000,
	}

	req5 := TestRequest{
		ClientID: 2412412412,
		Amount:   25000,
	}

	req6 := TestRequest{
		ClientID: 4444444444,
		Amount:   35000,
	}

	req7 := TestRequest{
		ClientID: 4444444444,
		Amount:   100000,
	}

	req8 := TestRequest{
		ClientID: 4444444444,
		Amount:   1000000000000,
	}
	AllReqs = append(AllReqs, req1, req2, req3, req4, req5, req6, req7, req8)

	for _, req := range AllReqs {
		reqBody, err := json.Marshal(req)
		if err != nil {
			t.Errorf("failed to marshal request %s", err.Error())
		}

		client := http.Client{}

		body := bytes.NewReader(reqBody)
		request, err := http.NewRequest("GET", "http://localhost:8080/create-transaction/get-money", body)
		if err != nil {
			t.Errorf("failed to create request %s", err.Error())
		}

		resp, err := client.Do(request)
		if err != nil {
			t.Errorf("faield to send request to server %s", err.Error())
		}

		buf := new(strings.Builder)
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			t.Errorf("failed to copy to buf %s", err.Error())
		}
		// check errors

		t.Logf("status code = %d", resp.StatusCode)
		t.Logf("result = %s", buf.String())
	}
}

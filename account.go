package jenga

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AccountBalanceResponse struct {
	Currency string     `json:"currency"`
	Balances []*Balance `json:"balances"`
}

type Balance struct {
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

func (j *Jenga) GetAccountBalance(ctx context.Context, countryCode string) (*AccountBalanceResponse, error) {
	j.CheckPreconditions()

	reqUrl := fmt.Sprintf("%s/v3-apis/account-api/v3.0/accounts/balances/%s/%s",
		JENGA_API_HOST, countryCode, j.JengaAuthCredentials.MerchantCode)

	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create get account balance request: %w", err)
	}

	signature, err := j.generateSignature(countryCode, j.JengaAuthCredentials.MerchantCode)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", j.JengaToken.AccessToken))
	req.Header.Set("Signature", signature)

	fmt.Println(signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	err = ErrorHandler(resp)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body could not be read: %w", err)
	}

	var accountBalanceResponse AccountBalanceResponse
	err = json.Unmarshal(body, &accountBalanceResponse)
	if err != nil {
		return nil, fmt.Errorf("response body could not be unmarshalled: %w", err)
	}

	return &accountBalanceResponse, nil

}

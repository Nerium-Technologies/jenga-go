package jenga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

const (
	JENGA_MERCHANT_CODE   = "JENGA_MERCHANT_CODE"
	JENGA_CONSUMER_SECRET = "JENGA_CONSUMER_SECRET"
	JENGA_API_KEY         = "JENGA_API_KEY"
	PRIVATE_KEY_PATH      = "PRIVATE_KEY_FILE_PATH"

	JENGA_API_HOST = "https://uat.finserve.africa"
)

type Jenga struct {
	JengaAuthCredentials *JengaAuthCredentials
	JengaToken           *JengaToken
}

type JengaAuthCredentials struct {
	MerchantCode       string `json:"merchantCode" validate:"required"`
	ConsumerSecret     string `json:"consumerSecret" validate:"required"`
	ApiKey             string `validate:"required"`
	PrivateKeyFilePath string `validate:"required"`
}

type JengaToken struct {
	AccessToken  string    `json:"accessToken" validate:"required"`
	RefreshToken string    `json:"refreshToken" validate:"required"`
	TokenType    string    `json:"tokenType" validate:"required"`
	IssuedAt     time.Time `json:"issuedAt" validate:"required"`
	ExpiresIn    time.Time `json:"expiresIn" validate:"required"`
}

type ValidatorService struct{}

func (vs *ValidatorService) Validate(s interface{}) error {
	v := validator.New()
	err := v.Struct(s)

	if err != nil {
		var sb strings.Builder
		for i, e := range err.(validator.ValidationErrors) {

			if i == len(err.(validator.ValidationErrors))-1 {
				sb.WriteString(fmt.Sprintf("%s is %s", e.Field(), e.Tag()))
			} else {
				sb.WriteString(fmt.Sprintf("%s is %s, ", e.Field(), e.Tag()))
			}

		}
		return fmt.Errorf("[%s]", sb.String())
	}

	return nil
}

func (j *Jenga) CheckPreconditions() {
	v := &ValidatorService{}

	// Check if credentials provided are valid
	err := v.Validate(j.JengaToken)
	if err != nil {
		log.Fatalf("invalid jenga auth token: %v", err)
	}

	// Check if credentials provided are valid
	err = v.Validate(j.JengaAuthCredentials)
	if err != nil {
		log.Fatalf("invalid jenga auth creds: %v", err)
	}
}

func New(creds *JengaAuthCredentials) (*Jenga, error) {
	var jenga Jenga

	v := &ValidatorService{}

	// Check if credentials provided are valid
	err := v.Validate(creds)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials provided: %w", err)
	}

	jengaToken, err := jenga.GenerateJengaBearerToken(creds)
	if err != nil {
		log.Fatalf("could not generate jenga bearer token: %v", err)
	}

	err = v.Validate(jengaToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token retrieved: %w", err)
	}

	jenga.JengaToken = jengaToken
	jenga.JengaAuthCredentials = creds

	return &jenga, nil
}

func (j *Jenga) GenerateJengaBearerToken(creds *JengaAuthCredentials) (*JengaToken, error) {

	jsonCreds, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("could not marshall credentials: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://uat.finserve.africa/authentication/api/v3/authenticate/merchant",
		bytes.NewBuffer(jsonCreds),
	)
	if err != nil {
		return nil, fmt.Errorf("new request could not be generated: %w", err)
	}

	req.Header.Set("Api-Key", fmt.Sprintf("%s", creds.ApiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body could not be read: %w", err)
	}

	err = ErrorHandler(resp)
	if err != nil {
		return nil, err
	}

	var jengaToken JengaToken
	err = json.Unmarshal(body, &jengaToken)
	if err != nil {
		return nil, fmt.Errorf("response body could not be unmarshalled: %w", err)
	}

	return &jengaToken, nil
}

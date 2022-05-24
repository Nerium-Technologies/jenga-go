package main

import (
	"context"
	"fmt"

	"github.com/Nerium-Technologies/jenga-go"
)

func main() { // nolint

	apiKey := jenga.MustGetEnvVar(jenga.JENGA_API_KEY)
	consumerSecret := jenga.MustGetEnvVar(jenga.JENGA_CONSUMER_SECRET)
	merchantCode := jenga.MustGetEnvVar(jenga.JENGA_MERCHANT_CODE)
	privateKeyFilePath := jenga.MustGetEnvVar(jenga.PRIVATE_KEY_PATH)

	j, err := jenga.New(&jenga.JengaAuthCredentials{
		ApiKey:             apiKey,
		ConsumerSecret:     consumerSecret,
		MerchantCode:       merchantCode,
		PrivateKeyFilePath: privateKeyFilePath,
	})

	if err != nil {
		fmt.Println("could not initialise jenga: ", err)
	}

	fmt.Printf("\ntoken: %v\n", j.JengaToken.AccessToken[:6])

	ctx := context.Background()

	bal, err := j.GetAccountBalance(ctx, "KE")
	if err != nil {
		fmt.Println("\ncould not get account balance: ", err)
	}

	fmt.Printf("\n%v\n", bal)
}

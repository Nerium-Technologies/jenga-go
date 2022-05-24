package jenga

import (
	"log"
	"os"
)

func MustGetEnvVar(key string) string {

	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("could not get envar of key: [%s]", key)
	}

	return value
}

package env

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string, defaultVal string) string {

	// load .env file
	err := godotenv.Load("env/.env")

	if err != nil {
		log.Println("Error loading .env file", err.Error())
	}

	val := os.Getenv(key)
	//log.Println("Got ", key, val)
	if val == "" {
		val = defaultVal
	}

	return val
}

func IsInDebugMode() bool {
	debugS := strings.TrimSpace(strings.ToUpper(GetEnvVariable("DEBUG", "FALSE")))

	if debugS == "TRUE" || debugS == "YES" || debugS == "Y" {
		return true
	}

	return false
}

func UseHttps() bool {
	debugS := strings.TrimSpace(strings.ToUpper(GetEnvVariable("HTTPS", "TRUE")))

	if debugS == "TRUE" || debugS == "YES" || debugS == "Y" {
		return true
	}

	return false
}

func GetSMTPServerData() (string, int, string, string) {
	host := strings.TrimSpace(GetEnvVariable("SMTP_HOST", ""))
	port := strings.TrimSpace(GetEnvVariable("SMTP_PORT", "587"))

	portX, err := strconv.Atoi(port)
	if err != nil || portX == 0 {
		portX = 587
	}

	user := strings.TrimSpace(GetEnvVariable("SMTP_USERNAME", ""))
	pwd := strings.TrimSpace(GetEnvVariable("SMTP_PASSWORD", ""))

	return host, portX, user, pwd

}

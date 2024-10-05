package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/brianvoe/gofakeit/v6/data"
)

func FakeHTTPStatusCode() string {
	return strconv.Itoa(gofakeit.HTTPStatusCode())
}

func FakeVISACard() string {

	// []string{"visa", "mastercard", "american-express", "diners-club", "discover", "jcb", "unionpay", "maestro", "elo", "hiper", "hipercard"}

	return gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: []string{"visa"}})

}
func FakeMasterCard() string {

	// []string{"visa", "mastercard", "american-express", "diners-club", "discover", "jcb", "unionpay", "maestro", "elo", "hiper", "hipercard"}

	return gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: []string{"mastercard"}})

}

func FakeCreditCard() string {

	// []string{"visa", "mastercard", "american-express", "diners-club", "discover", "jcb", "unionpay", "maestro", "elo", "hiper", "hipercard"}

	return gofakeit.CreditCardNumber(&gofakeit.CreditCardOptions{Types: data.CreditCardTypes})

}

var RandomFunctionMap map[string]func() string = map[string]func() string{
	"NAME":      gofakeit.Name,
	"FIRSTNAME": gofakeit.FirstName,
	"LASTNAME":  gofakeit.LastName,
	"USERNAME":  gofakeit.Username,

	"EMAIL":          gofakeit.Email,
	"URL":            gofakeit.URL,
	"DOMAIN":         gofakeit.DomainName,
	"IPV4":           gofakeit.IPv4Address,
	"IPV6":           gofakeit.IPv6Address,
	"HTTPSTATUSCODE": FakeHTTPStatusCode,

	"PHONE":   gofakeit.Phone,
	"CITY":    gofakeit.City,
	"STATE":   gofakeit.StateAbr,
	"ZIP":     gofakeit.Zip,
	"WORD":    gofakeit.Word,
	"SETENCE": gofakeit.SentenceSimple,

	"VISACARD":   FakeVISACard,
	"MASTERCARD": FakeMasterCard,
	"CREDITCARD": FakeCreditCard,
	"CARDCVV":    gofakeit.CreditCardCvv,
	"CARDEXPIRY": gofakeit.CreditCardExp,

	"BANKROUTING": gofakeit.AchRouting,
	"BANKACCOUNT": gofakeit.AchAccount,

	"COMPANYNAME": gofakeit.Company,

	"APPNAME": gofakeit.AppName,

	"COLOR": gofakeit.Color,

	"LANGUAGE": gofakeit.Language,

	"NUMBER": gofakeit.Digit,
}

func GetRandonFunctiolist() []string {
	returnMap := make([]string, 0)

	for k := range RandomFunctionMap {
		key := fmt.Sprintf("*RANDOM:%s", k)
		returnMap = append(returnMap, key)
	}

	return returnMap
}

func GenerateRandom(what string) (string, error) {
	funcToUse, found := RandomFunctionMap[strings.TrimSpace(what)]
	if found {
		return funcToUse(), nil
	}
	return "", fmt.Errorf("RANDOM Identifier not found:%s", what)
}

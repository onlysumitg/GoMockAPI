package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/onlysumitg/GoMockAPI/env"
)

// ------------------------------------------------------
//
// ------------------------------------------------------
type parameters struct {
	host           string
	port           int
	superuseremail string
	superuserpwd   string

	testmode bool
	domain   string

	useletsencrypt bool
	https          bool
	//staticDir string
	//flag      bool
}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (p *parameters) getHttpAddress() (string, string) {
	addr := p.host

	if p.port > 0 {
		addr = fmt.Sprintf("%s:%d", addr, p.port)
	}
	protocol := "http://"

	if p.https {
		protocol = "https://"

	}

	//if p.domain == "localhost" || p.domain == "0.0.0.0" {
	domain := p.domain
	if p.port > 0 {
		domain = fmt.Sprintf("%s:%d", p.domain, p.port)
	}

	return addr, fmt.Sprintf("%s%s", protocol, domain)
}

// ------------------------------------------------------
//
// ------------------------------------------------------
func (params *parameters) Load() {
	flag.StringVar(&params.host, "host", "", "Http Host Name")
	flag.IntVar(&params.port, "port", 4041, "Port")

	flag.StringVar(&params.superuseremail, "superuseremail", "admin2@example.com", "Super User email")
	flag.StringVar(&params.superuserpwd, "superuserpwd", "adminpass", "Super User password")

	flag.BoolVar(&params.useletsencrypt, "useletsencrypt", false, "Use let's encrypt ssl certificate")

	domain := "0.0.0.0"
	if runtime.GOOS == "windows" {
		domain = "localhost"
	}

	flag.StringVar(&params.domain, "domain", domain, "Domain name")

	flag.Parse()

	envPort := env.GetEnvVariable("PORT", "")

	port, err := strconv.Atoi(envPort)
	if err == nil {
		params.port = port
	}

	domainEnv := env.GetEnvVariable("DOMAIN", "")

	if domainEnv != "" {
		params.domain = domainEnv
	}

	useletsencrypt := strings.ToUpper(env.GetEnvVariable("USELETSENCRYPT", ""))

	if useletsencrypt != "" {
		if useletsencrypt == "TRUE" || useletsencrypt == "YES" || useletsencrypt == "Y" {
			params.useletsencrypt = true
		}

	}

	params.https = env.UseHttps()
	//fmt.Println("params.domain", params.domain)
	params.testmode = env.IsInDebugMode()

}

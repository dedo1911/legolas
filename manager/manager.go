package manager

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/dedo1911/legolas/storage"
	"github.com/dedo1911/legolas/users"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

type CertificateRequest struct {
	Email     string
	AuthEmail string
	AuthKey   string
	Domain    string
	IsStaging bool
}

func GetCertificate(request *CertificateRequest) *certificate.Resource {
	user, err := users.GetOrCreateUser(request.Email)
	if err != nil {
		log.Panicln(err)
	}

	config := lego.NewConfig(user)
	config.CADirURL = lego.LEDirectoryProduction
	if request.IsStaging {
		config.CADirURL = lego.LEDirectoryStaging
	}
	config.Certificate.KeyType = certcrypto.RSA4096

	client, err := lego.NewClient(config)
	if err != nil {
		log.Panicln(err)
	}
	if err := user.Register(client); err != nil {
		log.Panicln(err)
	}

	// Try to get exsisting certificate
	crt, err := storage.GetCertificate(request.Email, request.Domain)
	isNewDomain := false
	if err != nil {
		if !os.IsNotExist(err) {
			log.Panicln(err)
		}
		isNewDomain = true
	}

	if isNewDomain { // Generate new certificate
		cfConfig := cloudflare.NewDefaultConfig() // TODO: move into request
		cfConfig.AuthEmail = request.AuthEmail
		cfConfig.AuthKey = request.AuthKey
		cfProvider, err := cloudflare.NewDNSProviderConfig(cfConfig)
		if err != nil {
			log.Panicln(err)
		}

		err = client.Challenge.SetDNS01Provider(cfProvider)
		if err != nil {
			log.Panicln(err)
		}

		obtainRequest := certificate.ObtainRequest{
			Domains: []string{request.Domain},
			Bundle:  true,
		}

		certificates, err := client.Certificate.Obtain(obtainRequest)
		if err != nil {
			log.Panicln(err)
		}
		c, err := client.Certificate.Get(certificates.CertURL, true)
		if err != nil {
			log.Panicln(err)
		}
		c.PrivateKey = certificates.PrivateKey

		// Store certificate for later
		if err := storage.StoreCertificate(request.Email, request.Domain, c); err != nil {
			log.Panicln(err)
		}

		return c
	}

	// Check if current certificate is still valid
	cpb, _ := pem.Decode(crt.Certificate)
	xCert, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		log.Panicln(err)
	}

	_, err = xCert.Verify(x509.VerifyOptions{
		DNSName: request.Domain,
	})
	if err == nil { // Stored certificate is still valid
		return crt
	}

	log.Println("ERR:", err)

	// Renew certificate
	crt, err = client.Certificate.Renew(*crt, true, false, "")
	if err != nil {
		log.Panicln(err)
	}
	c, err := client.Certificate.Get(crt.CertURL, true)
	if err != nil {
		log.Panicln(err)
	}
	c.PrivateKey = crt.PrivateKey

	// Store certificate for later
	if err := storage.StoreCertificate(request.Email, request.Domain, c); err != nil {
		log.Println(err)
	}

	return c
}

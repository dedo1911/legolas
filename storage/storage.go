package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certificate"
)

const certificatesFolder = "certificates"

type CertificateResource struct {
	Domain            string `json:"domain"`
	CertURL           string `json:"certUrl"`
	CertStableURL     string `json:"certStableUrl"`
	PrivateKey        []byte `json:"privateKey"`
	Certificate       []byte `json:"certificate"`
	IssuerCertificate []byte `json:"issuerCertificate"`
	CSR               []byte `json:"csr"`
}

func getFileName(account, domain string) string {
	hasher1 := sha1.New()
	hasher1.Write([]byte(account))
	hasher2 := sha1.New()
	hasher2.Write([]byte(domain))
	return filepath.Join(certificatesFolder, hex.EncodeToString(hasher1.Sum(nil)[:8])+"-"+hex.EncodeToString(hasher2.Sum(nil))+".json")
}

func GetCertificate(account, domain string) (*certificate.Resource, error) {
	os.MkdirAll(certificatesFolder, os.ModePerm)
	data, err := os.ReadFile(getFileName(account, domain))
	if err != nil {
		return nil, err
	}
	var resource CertificateResource
	if err := json.Unmarshal(data, &resource); err != nil {
		return nil, err
	}
	crt := certificate.Resource{
		Domain:            resource.Domain,
		CertURL:           resource.CertURL,
		CertStableURL:     resource.CertStableURL,
		PrivateKey:        resource.PrivateKey,
		Certificate:       resource.Certificate,
		IssuerCertificate: resource.IssuerCertificate,
		CSR:               resource.CSR,
	}
	return &crt, nil
}

func StoreCertificate(account, domain string, crt *certificate.Resource) error {
	resource := CertificateResource{
		Domain:            crt.Domain,
		CertURL:           crt.CertURL,
		CertStableURL:     crt.CertStableURL,
		PrivateKey:        crt.PrivateKey,
		Certificate:       crt.Certificate,
		IssuerCertificate: crt.IssuerCertificate,
		CSR:               crt.CSR,
	}
	data, err := json.MarshalIndent(resource, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(getFileName(account, domain), data, 0600)
}

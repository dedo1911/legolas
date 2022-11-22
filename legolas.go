package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dedo1911/legolas/manager"
	"github.com/dedo1911/legolas/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	fmt.Println("  _           _____       _            _____ ")
	fmt.Println(" | |         / ____|     | |          / ____|")
	fmt.Println(" | |     ___| |  __  ___ | |     __ _| (___  ")
	fmt.Println(" | |    / _ \\ | |_ |/ _ \\| |    / _` |\\___ \\ ")
	fmt.Println(" | |___|  __/ |__| | (_) | |___| (_| |____) |")
	fmt.Println(" |______\\___|\\_____|\\___/|______\\__,_|_____/ ")
	fmt.Println()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: move to template
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
		<html>
			<head>
				<title>LeGoLaS</title>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width,initial-scale=1"/>
				<style type="text/css">
				* {
					font-family: 'Lato', sans-serif;
				}
				body {
					margin: 0;
					padding: 0;
					height: 100vh;
					display: grid; 
					grid-template-columns: auto; 
					grid-template-rows: 8em 1fr 4em; 
					gap: 0px 0px; 
					grid-template-areas: 
						"header"
						"main"
						"footer";
				}
				header {
					grid-area: header;
					padding: 0 .5em;
					text-align: center;
					
					background: #FFBA08;
					color: #D00000;
				}
				main {
					grid-area: main;
					padding: 1em;
					text-align: center;

					font-size: large;
					background: #A2AEBB;
					color: #1C3144;
				}
				footer {	
					grid-area: footer;
					padding: 0 1em;
					text-align: right;

					background: #1C3144;
					color: #A2AEBB;
				}
				</style>
				<link rel="preconnect" href="https://fonts.googleapis.com">
				<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
				<link href="https://fonts.googleapis.com/css2?family=Lato:wght@300;700&display=swap" rel="stylesheet">
			</head>
			<body>
				<header>
					<h1>LeGoLaS: LeGo Listen and Serve</h1>
				</header>
				<main>
					<p>GET /certificate?email=<strong>{email}</strong>&authEmail=<strong>{authEmail}</strong>&authKey=<strong>{authKey}</strong>&domain=<strong>{domain}</strong></p>
				</main>
				<footer>
						<p>...and you have my bow!</p>
				</footer>
			</body>
		</html>`))
	})
	r.Get("/certificate", getCertificate)
	// TODO: make port configurable
	log.Println("Listening on port 3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Panicln(err)
	}
	// TODO: add graceful shutdown
}

func getCertificate(w http.ResponseWriter, r *http.Request) {
	// TODO: create apikey auth mechanism
	certificates, err := manager.GetCertificate(&manager.CertificateRequest{
		Email:     r.URL.Query().Get("email"),
		AuthEmail: r.URL.Query().Get("authEmail"),
		AuthKey:   r.URL.Query().Get("authKey"),
		Domain:    r.URL.Query().Get("domain"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, storage.CertificateResource{
		IssuerCertificate: certificates.IssuerCertificate,
		Certificate:       certificates.Certificate,
		PrivateKey:        certificates.PrivateKey,
	})
}

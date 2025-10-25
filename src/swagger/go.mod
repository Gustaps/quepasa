module github.com/nocodeleaks/quepasa/swagger

replace github.com/nocodeleaks/quepasa/api => ../api

replace github.com/nocodeleaks/quepasa/form => ../form

replace github.com/nocodeleaks/quepasa/models => ../models

replace github.com/nocodeleaks/quepasa/environment => ../environment

replace github.com/nocodeleaks/quepasa/library => ../library

replace github.com/nocodeleaks/quepasa/media => ../media

replace github.com/nocodeleaks/quepasa/metrics => ../metrics

replace github.com/nocodeleaks/quepasa/rabbitmq => ../rabbitmq

replace github.com/nocodeleaks/quepasa/webserver => ../webserver

replace github.com/nocodeleaks/quepasa/whatsapp => ../whatsapp

replace github.com/nocodeleaks/quepasa/whatsmeow => ../whatsmeow

replace github.com/nocodeleaks/quepasa/signalr => ../signalr

replace github.com/nocodeleaks/quepasa/sipproxy => ../sipproxy

replace github.com/nocodeleaks/quepasa/swagger => ../swagger

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/nocodeleaks/quepasa/environment v0.0.0-00010101000000-000000000000
	github.com/nocodeleaks/quepasa/webserver v0.0.0-00010101000000-000000000000
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.4
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/nocodeleaks/quepasa/library v0.0.0-00010101000000-000000000000 // indirect
	github.com/nocodeleaks/quepasa/whatsapp v0.0.0-00010101000000-000000000000 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

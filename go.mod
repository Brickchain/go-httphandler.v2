module github.com/IpsoVeritas/httphandler

go 1.14

replace github.com/IpsoVeritas/document => ../document

replace github.com/IpsoVeritas/crypto => ../crypto

replace github.com/IpsoVeritas/logger => ../logger

require (
	github.com/IpsoVeritas/crypto v0.0.0-20181010203950-c229a2b23e68
	github.com/IpsoVeritas/document v0.0.0-20180814075806-099bc71d4b53
	github.com/IpsoVeritas/logger v0.0.0-20180912100710-b76d97958f28
	github.com/gorilla/handlers v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	golang.org/x/text v0.3.6
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/square/go-jose.v1 v1.1.2
)

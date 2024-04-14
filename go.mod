module sjones/wyze-go

go 1.18

require github.com/go-resty/resty/v2 v2.10.0

require github.com/golang-jwt/jwt/v5 v5.2.1 // indirect

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/slack-go/slack v0.12.3
	golang.org/x/net v0.17.0 // indirect
	stevejones.softwaredev/wyze-go/wyze v1.0.0
)

replace stevejones.softwaredev/wyze-go/wyze v1.0.0 => ./wyze

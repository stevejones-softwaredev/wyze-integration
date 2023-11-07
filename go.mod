module sjones/wyze-go

go 1.18

require github.com/go-resty/resty/v2 v2.10.0

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.17.0 // indirect
	github.com/slack-go/slack v0.12.3
	sjones/wyze-go/wyze v1.0.0
)

replace sjones/wyze-go/wyze v1.0.0 => ./wyze

module github.com/Dreamacro/clash

go 1.17

require (
	github.com/Dreamacro/go-shadowsocks2 v0.1.7
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.0
	github.com/go-chi/render v1.0.1
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/insomniacslk/dhcp v0.0.0-20211214070828-5297eed8f489
	github.com/kentik/patricia v0.0.0-20201202224819-f9447a6e25f1
	github.com/miekg/dns v1.1.45
	github.com/oschwald/geoip2-golang v1.5.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.etcd.io/bbolt v1.3.6
	go.uber.org/atomic v1.9.0
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/net v0.0.0-20211216030914-fe4d6282115f
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e
	golang.zx2c4.com/wireguard/windows v0.5.1
	gopkg.in/yaml.v2 v2.4.0
	gvisor.dev/gvisor v0.0.0-20211214225704-fb4fed4724e4
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/oschwald/maxminddb-golang v1.8.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/u-root/uio v0.0.0-20210528114334-82958018845c // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/text v0.3.8-0.20211004125949-5bd84dd9b33b // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.1.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/Dreamacro/go-shadowsocks2 v0.1.7 => github.com/wwqgtxx/go-shadowsocks2 v0.1.7

module github.com/axli-personal/drive/backend/drive

go 1.19

require (
	github.com/axli-personal/drive/backend/common v0.0.0
	github.com/axli-personal/drive/backend/pkg v0.0.0
	github.com/caarlos0/env/v7 v7.0.0
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/google/uuid v1.3.0
	github.com/redis/go-redis/v9 v9.0.2
	github.com/sirupsen/logrus v1.9.0
	gorm.io/driver/mysql v1.4.7
	gorm.io/gorm v1.25.0
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)

replace (
	github.com/axli-personal/drive/backend/common => ../common/
	github.com/axli-personal/drive/backend/pkg => ../pkg/
)

module github.com/axli-personal/drive/backend/channel

go 1.19

require (
	github.com/axli-personal/drive/backend/common v0.0.0
    github.com/axli-personal/drive/backend/pkg v0.0.0
    github.com/google/uuid v1.3.0
)

replace (
	github.com/axli-personal/drive/backend/common => ../common/
	github.com/axli-personal/drive/backend/pkg => ../pkg/
)
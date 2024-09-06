module github.com/Nicknamezz00/gorder/internal/orders

go 1.22.5

replace github.com/Nicknamezz00/gorder/internal/common => ../common/

require (
	github.com/Nicknamezz00/gorder/internal/common v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)

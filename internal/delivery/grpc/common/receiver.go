package common

import "google.golang.org/grpc"

type Receiver interface {
	Register(server grpc.Server)
}

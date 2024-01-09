package sender

import (
	"context"
)

func SenderInstance(ctx context.Context) Sender {
	return ctx.Value(ServiceName).(*sender)
}

func CtxSenderInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceName, entry)
}

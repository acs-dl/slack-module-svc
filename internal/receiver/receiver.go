package receiver

import (
	"context"
)

func ReceiverInstance(ctx context.Context) Receiver {
	return ctx.Value(ServiceName).(*receiver)
}

func CtxReceiverInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceName, entry)
}

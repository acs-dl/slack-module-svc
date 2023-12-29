package worker

import "context"

func WorkerInstance(ctx context.Context) Worker {
	return ctx.Value(ServiceName).(*worker)
}

func CtxWorkerInstance(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, ServiceName, entry)
}

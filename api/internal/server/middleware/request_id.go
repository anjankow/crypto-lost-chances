package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type keyReqIDType string

const keyRequestID keyReqIDType = "requestID"

func addContextRequestID(ctx context.Context) context.Context {
	reqID := uuid.New()

	return context.WithValue(ctx, keyRequestID, reqID.String())
}

func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(keyRequestID)
	if id, ok := reqID.(string); ok {
		return id
	}

	return ""
}

func AddRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := addContextRequestID(r.Context())

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

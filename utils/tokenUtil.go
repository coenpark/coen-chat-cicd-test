package utils

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

func ExtractTokenFromContext(ctx context.Context) (string, string) {
	return metautils.ExtractIncoming(ctx).Get("access-token"), metautils.ExtractIncoming(ctx).Get("refresh-token")
}

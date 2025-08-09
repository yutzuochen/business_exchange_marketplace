package graphql

import (
	"context"
	"strings"

	"trade_company/internal/auth"
	"trade_company/internal/config"
)

type ctxKey string

const ctxUserIDKey ctxKey = "graphqlUserID"

func WithUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (uint, bool) {
	v := ctx.Value(ctxUserIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

// ExtractUserFromAuthHeader parses Authorization header and embeds user ID to ctx if valid.
func ExtractUserFromAuthHeader(cfg *config.Config, parent context.Context, authorizationHeader string) context.Context {
	if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return parent
	}
	token := strings.TrimPrefix(authorizationHeader, "Bearer ")
	claims, err := auth.ParseToken(cfg, token)
	if err != nil {
		return parent
	}
	return WithUserID(parent, claims.UserID)
}

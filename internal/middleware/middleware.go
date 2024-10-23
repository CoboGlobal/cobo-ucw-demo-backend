package middleware

import (
	"cobo-ucw-backend/internal/middleware/auth"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(auth.NewJwtMiddleware)

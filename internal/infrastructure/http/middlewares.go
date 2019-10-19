package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/lzakharov/goss/internal/domain"
	routing "github.com/qiangxue/fasthttp-routing"
	"go.uber.org/zap"
)

const mimeJSON = "application/json"

type getClaims func(accessToken string) (*domain.AccessTokenClaims, error)

func jsonWriterMiddleware(ctx *routing.Context) error {
	ctx.SetContentType(mimeJSON)
	ctx.Serialize = json.Marshal
	return nil
}

func loggerMiddleware(logger *zap.Logger) routing.Handler {
	return func(ctx *routing.Context) error {
		uri := ctx.Request.URI().String()
		method := string(ctx.Method())

		u, err := uuid.NewRandom()
		if err != nil {
			logger.Error("Error generating a request id!",
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Error(err))
			return err
		}
		requestID := u.String()

		ctx.Set(ctxRequestID, requestID)

		logger.Info("Got request.",
			zap.String("requestID", requestID),
			zap.String("uri", uri),
			zap.String("method", method))

		if err := ctx.Next(); err != nil {
			logger.Error("Error handling the request!",
				zap.String("requestID", requestID),
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Error(err))
			return err
		}

		return nil
	}
}

func authMiddleware(getClaims getClaims) routing.Handler {
	return func(ctx *routing.Context) error {
		accessToken := string(ctx.Request.Header.Peek(authorizationHeader))

		claims, err := getClaims(accessToken)
		if err != nil {
			ctx.SetStatusCode(http.StatusUnauthorized)
			return ctx.WriteData(err)
		}

		ctx.Set(ctxClaims, claims)
		return nil
	}
}

func errorHandlerMiddleware(ctx *routing.Context) error {
	if err := ctx.Next(); err != nil {
		requestID := ctx.Get(ctxRequestID).(string)

		switch err {
		case domain.ErrInternalStorage:
			ctx.SetStatusCode(http.StatusInternalServerError)
		case domain.ErrNotFound:
			ctx.SetStatusCode(http.StatusNotFound)
		case domain.ErrInvalidCredentials:
			ctx.SetStatusCode(http.StatusUnauthorized)
		case domain.ErrInternalSecurity:
			ctx.SetStatusCode(http.StatusInternalServerError)
		case domain.ErrInvalidAccessToken:
			ctx.SetStatusCode(http.StatusUnauthorized)
		case domain.ErrInvalidRefreshToken:
			ctx.SetStatusCode(http.StatusBadRequest)
		default:
			resp := &ErrorResponse{
				Status:    http.StatusInternalServerError,
				Message:   err.Error(),
				RequestID: requestID,
			}

			return ctx.WriteData(resp)
		}

		resp := &ErrorResponse{
			Status:    errStatus[err],
			Message:   err.Error(),
			RequestID: requestID,
		}

		return ctx.WriteData(resp)
	}

	return nil
}

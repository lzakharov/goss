package http

import (
	"encoding/json"
	"net/http"

	"github.com/lzakharov/goss/internal/domain"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func (a *adapter) newRouter() fasthttp.RequestHandler {
	router := routing.New()
	router.Use(loggerMiddleware(a.logger), jsonWriterMiddleware, errorHandlerMiddleware)

	authMiddleware := authMiddleware(a.service.GetAccessTokenClaims)

	v1 := router.Group("/v1")
	{
		v1.Get("/health", a.Health)

		auth := v1.Group("/auth")
		{
			auth.Post("/login", a.Login)
			auth.Post("/refresh", a.Refresh)
		}

		user := v1.Group("/user")
		user.Use(authMiddleware)
		{
			user.Get("/self", a.GetUser)
			user.Post("/logout", a.Logout)
		}
	}

	return router.HandleRequest
}

//Health responses with the service health status.
func (a *adapter) Health(ctx *routing.Context) error {
	health := a.service.CheckHealth()
	return ctx.WriteData(health)
}

//Login handles user login.
func (a *adapter) Login(ctx *routing.Context) error {
	var credentials *domain.Credentials

	if err := json.Unmarshal(ctx.Request.Body(), &credentials); err != nil {
		a.logger.Error("Error unmarshalling credentials!", zap.Binary("body", ctx.Request.Body()), zap.Error(err))
		return err
	}

	if err := a.validator.Struct(credentials); err != nil {
		a.logger.Error("Invalid credentials!", zap.Any("credentials", credentials), zap.Error(err))
		return err
	}

	authData, err := a.service.Login(credentials)
	if err != nil {
		a.logger.Error("Login error!", zap.Error(err))
		return err
	}

	return ctx.WriteData(authData)
}

//Refresh handles user access token refresh.
func (a *adapter) Refresh(ctx *routing.Context) error {
	var refreshToken string

	if err := json.Unmarshal(ctx.Request.Body(), &refreshToken); err != nil {
		a.logger.Error("Error unmarshalling a refresh token!", zap.String("refresh", refreshToken), zap.Error(err))
		return err
	}

	authData, err := a.service.RefreshToken(refreshToken)
	if err != nil {
		a.logger.Error("Error refreshing a refresh token!", zap.String("refresh", refreshToken), zap.Error(err))
		return err
	}

	return ctx.WriteData(authData)
}

//GetUser returns current logged in user.
func (a *adapter) GetUser(ctx *routing.Context) error {
	claims := ctx.Get(ctxClaims).(*domain.AccessTokenClaims)

	user, err := a.service.GetUser(claims.UserID)
	if err != nil {
		a.logger.Error("Error getting the logged in user!",
			zap.Any("claims", claims),
			zap.Error(err))
		return err
	}

	return ctx.WriteData(user)
}

//Logout handles user logout.
func (a *adapter) Logout(ctx *routing.Context) error {
	claims := ctx.Get(ctxClaims).(*domain.AccessTokenClaims)

	if err := a.service.Logout(claims.UserID); err != nil {
		a.logger.Error("Logout error!",
			zap.Any("claims", claims),
			zap.Error(err))
		return err
	}

	ctx.SetStatusCode(http.StatusNoContent)
	return nil
}

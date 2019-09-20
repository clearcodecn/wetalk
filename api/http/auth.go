package http

import (
	"context"
	"github.com/clearcodecn/wetalk/api/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) authMiddleware(ctx *gin.Context) {
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, s.keyFunc)
	if err != nil {
		ctx.JSON(401, fail("unauthorized", nil))
		return
	}
	if !token.Valid {
		ctx.JSON(401, fail("token is not valid", nil))
		return
	}
	cc := context.Background()
	cc = context.WithValue(cc, TokenCtx, token)
	ctx.Request = ctx.Request.WithContext(cc)
}

func (s *Server) keyFunc(token *jwt.Token) (i interface{}, e error) {
	return []byte(s.config.HttpConfig.JwtKey), nil
}

var (
	GetUserCtx = struct{}{}
	TokenCtx   = struct{}{}
)

type GetUserFunc func(context context.Context)

// GetUser is a lazy func to get user from context
func (s *Server) GetUser(ctx context.Context) (*model.User, context.Context, error) {
	user, _ := ctx.Value(GetUserCtx).(*model.User)
	if user != nil {
		return user, ctx, nil
	}
	token, ok := ctx.Value(TokenCtx).(*jwt.Token)
	if !ok {
		return nil, ctx, errors.New(`no token in context`)
	}
	id := token.Claims.(jwt.MapClaims)["id"].(int)
	user, err := s.model.GetUserById(id)
	if err != nil {
		return nil, ctx, err
	}
	ctx = context.WithValue(ctx, GetUserCtx, user)
	return user, ctx, nil
}

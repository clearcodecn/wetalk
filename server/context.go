package server

import "golang.org/x/net/context"

type ConnectionContext struct {
	context.Context
}

func (ctx *ConnectionContext) Close() error {
	
}
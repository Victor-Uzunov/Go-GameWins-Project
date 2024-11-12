package graph

import (
	"application/graph/model"
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
)

func RoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, requiredRole model.Role) (res interface{}, err error) {
	userRole, ok := ctx.Value("role").(model.Role)
	if !ok {
		return nil, errors.New("no role found in context")
	}

	if userRole != requiredRole {
		return nil, fmt.Errorf("access denied, requires %s role", requiredRole)
	}

	return next(ctx)
}

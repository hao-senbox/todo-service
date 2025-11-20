package repair

import (
	"context"
	"fmt"
	"strings"
	"todo-service/pkg/constants"

	"github.com/golang-jwt/jwt/v5"
)

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) CanUpdateReport(ctx context.Context, repair *Repair, userID string) error {
	if repair.ReportBy != userID {
		return fmt.Errorf("user does not have permission to update this report")
	}
	return nil
}

func (p *Policy) hasAdminRole(tokenString string) bool {
	token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if roles, ok := claims["roles"].(string); ok {
			roleList := strings.Split(roles, ", ")
			for _, role := range roleList {
				if strings.TrimSpace(role) == "Admin" {
					return true
				}
			}
		}
	}
	return false
}

func (p *Policy) CanAssignRepair(ctx context.Context, repair *Repair, userID string) error {
	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	if !p.hasAdminRole(token) {
		return fmt.Errorf("user does not have admin permission to assign repairs")
	}

	return nil
}

func (p *Policy) CanDeleteRepair(ctx context.Context, repair *Repair, userID string) error {
	if repair.ReportBy != userID {
		return fmt.Errorf("user does not have permission to delete this repair")
	}
	return nil
}

func (p *Policy) CanCompleteRepair(ctx context.Context, repair *Repair, userID string) error {
	if repair.AssignedTo != nil && *repair.AssignedTo != userID {
		return fmt.Errorf("user does not have permission to complete this repair")
	}
	return nil
}

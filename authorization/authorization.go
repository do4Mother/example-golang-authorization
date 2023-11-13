package authorization

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Action string
type Map map[string]interface{}
type Hook func() Map
type Resource func(c *fiber.Ctx) (Map, error)

const (
	Create Action = "create"
	Read   Action = "read"
	Update Action = "update"
	Delete Action = "delete"
)

type RBAC struct {
	Role string
	Action
	Object string
}

type Authorization struct {
	RBAC []RBAC
}

func (a *Authorization) Add(role string, action Action, object string) {
	a.RBAC = append(a.RBAC, RBAC{
		Role:   role,
		Action: action,
		Object: object,
	})
}

func (a *Authorization) IsAuthorized(role string, action Action, resource Map, data Map) bool {
	for _, rbac := range a.RBAC {
		if rbac.Role == role && rbac.Action == action {
			if rbac.Object == "*" {
				return true
			}

			// compare resource with data from rbac.Object as key of map
			if resource != nil && data != nil {
				for key, value := range resource {
					if rbac.Object == key {
						if value == data[key] {
							return true
						}
					}
				}
			}

		}
	}
	return false
}

func (a *Authorization) Middleware(action Action, resource Resource, hook Hook) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		var (
			source Map
			data   Map
		)

		if resource != nil {
			var err error
			source, err = resource(c)

			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
		}

		if hook != nil {
			data = hook()
		}

		if a.IsAuthorized(claims["role"].(string), action, source, data) {
			return c.Next()
		}

		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}
}

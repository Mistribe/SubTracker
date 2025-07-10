package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
)

type FamilyMemberDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteFamilyMemberCommand, bool]
}

func (f FamilyMemberDeleteEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid id format",
		})
		return
	}

	cmd := command.DeleteFamilyMemberCommand{
		Id: id,
	}

	r := f.handler.Handle(c, cmd)
	handleResponse(c, r, withMapping[bool](func(success bool) any {
		return map[string]bool{"success": success}
	}))
}

func (f FamilyMemberDeleteEndpoint) Pattern() []string {
	return []string{
		"/members/:id",
	}
}

func (f FamilyMemberDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (f FamilyMemberDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberDeleteEndpoint(handler core.CommandHandler[command.DeleteFamilyMemberCommand, bool]) *FamilyMemberDeleteEndpoint {
	return &FamilyMemberDeleteEndpoint{
		handler: handler,
	}
}

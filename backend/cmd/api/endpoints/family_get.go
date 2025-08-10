package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, family.Family]
}

// Handle godoc
//
//	@Summary		Get family by ID
//	@Description	Retrieve a family and its members by family ID
//	@Tags			family
//	@Produce		json
//	@Param			familyId	path		string				true	"Family ID (UUID format)"
//	@Success		200			{object}	familyModel			"Successfully retrieved family"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		401			{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse	"Family not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families/{familyId} [get]
func (f FamilyGetEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("familyId")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: "invalid id format",
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	q := query.FindOneQuery{
		Id: id,
	}

	r := f.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[family.Family](func(fm family.Family) any {
			return newFamilyModel(userId, fm)
		}))
}

func (f FamilyGetEndpoint) Pattern() []string {
	return []string{
		"/:familyId",
	}
}

func (f FamilyGetEndpoint) Method() string {
	return http.MethodGet
}

func (f FamilyGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberGetEndpoint(handler core.QueryHandler[query.FindOneQuery, family.Family]) *FamilyGetEndpoint {
	return &FamilyGetEndpoint{
		handler: handler,
	}
}

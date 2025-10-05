package currency

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/currency"

	dcurrency "github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type SupportedEndpoint struct {
}

func NewSupportedEndpoint() *SupportedEndpoint {
	return &SupportedEndpoint{}
}

// Handle godoc
//
//	@Summary		Get Supported Currencies
//	@Description	get details of all supported currencies
//	@Tags			currencies
//	@Produce		json
//	@Success		200	{array}	string	"currencies"
//	@Router			/currencies/supported [get]
func (e SupportedEndpoint) Handle(c *gin.Context) {
	response := herd.Select(dcurrency.GetSupportedCurrencies(), func(u currency.Unit) string {
		return u.String()
	})
	c.JSON(http.StatusOK, response)
}

func (e SupportedEndpoint) Pattern() []string {
	return []string{
		"/supported",
	}
}

func (e SupportedEndpoint) Method() string {
	return http.MethodGet
}

func (e SupportedEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

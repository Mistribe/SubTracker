package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dcurrency "github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"golang.org/x/text/currency"
)

type CurrencySupportedEndpoint struct {
}

func NewCurrencySupportedEndpoint() *CurrencySupportedEndpoint {
	return &CurrencySupportedEndpoint{}
}

// CurrencySupportedEndpoint godoc
//
//	@Summary		Get Supported Currencies
//	@Description	get details of all supported currencies
//	@Produce		json
//	@Success		200	{array}	string	"currencies"
//	@Router			/currencies/supported [get]
func (e CurrencySupportedEndpoint) Handle(c *gin.Context) {
	response := slicesx.Select(dcurrency.GetSupportedCurrencies(), func(u currency.Unit) string {
		return u.String()
	})
	c.JSON(http.StatusOK, response)
}

func (e CurrencySupportedEndpoint) Pattern() []string {
	return []string{
		"/supported",
	}
}

func (e CurrencySupportedEndpoint) Method() string {
	return http.MethodGet
}

func (e CurrencySupportedEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

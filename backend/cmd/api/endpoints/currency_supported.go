package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oleexo/subtracker/internal/domain/currency"
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
	c.JSON(http.StatusOK, currency.GetSupportedCurrencies())
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

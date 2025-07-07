package endpoints

import (
    "github.com/gin-gonic/gin"
    "github.com/oleexo/subtracker/internal/domain/subscription"
    "github.com/oleexo/subtracker/internal/infrastructure/api/ginfx"
)

type SubscriptionEndpointGroup struct {
    routes []ginfx.Route
}

func (s SubscriptionEndpointGroup) Prefix() string {
    return "/subscription"
}

func (s SubscriptionEndpointGroup) Routes() []ginfx.Route {
    return s.routes
}

func (s SubscriptionEndpointGroup) Middlewares() []gin.HandlerFunc {
    return nil
}

func NewSubscriptionEndpointGroup(getEndpoint *SubscriptionGetEndpoint) *SubscriptionEndpointGroup {
    return &SubscriptionEndpointGroup{
        routes: []ginfx.Route{
            getEndpoint,
        },
    }
}

type subscriptionModel struct {
    Id string `json:id`
}

func newSubscriptionModel(source subscription.Subscription) subscriptionModel {
    return subscriptionModel{
        Id: source.Id(),
    }
}

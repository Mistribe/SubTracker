package authentication

import (
	"context"

	"github.com/Oleexo/config-go"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x"
)

type clerkIdentityProvider struct {
	jwksClient *jwks.Client
	userClient *user.Client
}

type publicUserMetadata struct {
	Role string `json:"role"`
}

func NewClerkIdentityProvider(cfg config.Configuration) ports.IdentityProvider {
	secret := cfg.GetString("CLERK_SECRET_KEY")
	clerkConfig := &clerk.ClientConfig{}
	clerkConfig.Key = x.P(secret)
	return &clerkIdentityProvider{
		jwksClient: jwks.NewClient(clerkConfig),
		userClient: user.NewClient(clerkConfig),
	}
}

func (c *clerkIdentityProvider) DeleteUser(ctx context.Context, userId types.UserID) error {
	_, err := c.userClient.Delete(ctx, userId.String())
	return err
}

func (c *clerkIdentityProvider) ReadSessionToken(ctx context.Context, sessionToken string) (ports.Identity, error) {
	// Decode the session JWT to find the key LabelID
	unsafeClaims, err := jwt.Decode(ctx, &jwt.DecodeParams{
		Token: sessionToken,
	})
	if err != nil {
		return ports.NewInvalidIdentity(), err
	}

	// Fetch the JSON Web Key
	jwk, err := jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
		KeyID:      unsafeClaims.KeyID,
		JWKSClient: c.jwksClient,
	})
	if err != nil {
		return ports.NewInvalidIdentity(), authorization.ErrUnauthorized
	}

	// Verify the session
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: sessionToken,
		JWK:   jwk,
	})
	if err != nil {
		return ports.NewInvalidIdentity(), authorization.ErrUnauthorized
	}

	// todo extract metadata from public metadata
	return ports.Identity{
		Id:      claims.Subject,
		IsValid: true,
	}, nil
}

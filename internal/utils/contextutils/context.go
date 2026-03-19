package contextutils

import (
	"context"
	"errors"

	"evo-ai-core-service/internal/types"

	"github.com/google/uuid"
)

var ErrUnauthorized = errors.New("unauthorized")

func GetAccountID(ctx context.Context) (uuid.UUID, error) {
	if accountID, ok := ctx.Value("account_id").(uuid.UUID); ok {
		return accountID, nil
	}

	return uuid.Nil, errors.New("account_id not found in context")
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if userID, ok := ctx.Value("user_id").(uuid.UUID); ok {
		return userID, nil
	}

	return uuid.Nil, errors.New("user_id not found in context")
}

func GetUserEmail(ctx context.Context) (string, error) {
	if email, ok := ctx.Value("email").(string); ok {
		return email, nil
	}

	return "", errors.New("email not found in context")
}

func GetUserName(ctx context.Context) (string, error) {
	if name, ok := ctx.Value("name").(string); ok {
		return name, nil
	}

	return "", errors.New("name not found in context")
}

func GetToken(ctx context.Context) (string, error) {
	if token, ok := ctx.Value("token").(string); ok {
		return token, nil
	}

	return "", errors.New("token not found in context")
}

func GetApiAccessToken(ctx context.Context) (string, error) {
	if token, ok := ctx.Value("api_access_token").(string); ok {
		return token, nil
	}

	return "", errors.New("token not found in context")
}

func GetTokenType(ctx context.Context) (string, error) {
	if tokenType, ok := ctx.Value("token_type").(string); ok {
		return tokenType, nil
	}

	return "", errors.New("token type not found in context")
}

func GetAuthHeaders(ctx context.Context) (interface{}, error) {
	if headers := ctx.Value("auth_headers"); headers != nil {
		return headers, nil
	}

	return nil, errors.New("auth_headers not found in context")
}

// GetUser returns the complete user information from context
func GetUser(ctx context.Context) (types.EvoAuthUser, error) {
	if user, ok := ctx.Value("user").(types.EvoAuthUser); ok {
		return user, nil
	}

	return types.EvoAuthUser{}, errors.New("user not found in context")
}

// GetAccounts returns all accounts from context
func GetAccounts(ctx context.Context) ([]types.EvoAuthAccount, error) {
	if accounts, ok := ctx.Value("accounts").([]types.EvoAuthAccount); ok {
		return accounts, nil
	}

	return nil, errors.New("accounts not found in context")
}

// GetAccount returns the currently active account based on account_id in context
func GetAccount(ctx context.Context) (types.EvoAuthAccount, error) {
	if account, ok := ctx.Value("account").(types.EvoAuthAccount); ok {
		return account, nil
	}

	return types.EvoAuthAccount{}, errors.New("active account not found")
}


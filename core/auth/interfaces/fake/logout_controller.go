package fake

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/auth/application/fake"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//LogoutController fake implementation
	LogoutController struct {
		responder      *web.Responder
		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
	}
)

// Inject dependencies
func (l *LogoutController) Inject(
	responder *web.Responder,
	authManager *application.AuthManager,
	eventPublisher *application.EventPublisher,
) {
	l.responder = responder
	l.authManager = authManager
	l.eventPublisher = eventPublisher
}

// Get HTTP action
func (l *LogoutController) Get(ctx context.Context, request *web.Request) web.Result {
	request.Session().Delete(fake.UserSessionKey)
	l.eventPublisher.PublishLogoutEvent(ctx, &domain.LogoutEvent{
		Session: request.Session(),
	})

	redirectURL, _ := l.authManager.URL(ctx, "")

	return l.responder.URLRedirect(redirectURL)
}
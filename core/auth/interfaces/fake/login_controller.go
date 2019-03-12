package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// LoginController fake implementation
	LoginController struct {
		responder     *web.Responder
		authManager   *application.AuthManager
		loginTemplate string
	}
)

// Inject dependencies
func (l *LoginController) Inject(
	responder *web.Responder,
	authManager *application.AuthManager,
	cfg *struct {
		FakeLoginTemplate string `inject:"config:auth.fakeLoginTemplate"`
	},
) {
	l.responder = responder
	l.authManager = authManager
	l.loginTemplate = cfg.FakeLoginTemplate
}

// Get http action
func (l *LoginController) Get(ctx context.Context, request *web.Request) web.Result {
	redirectURL, ok := request.Params["redirecturl"]
	if !ok || redirectURL == "" {
		redirectURL = request.Request().Referer()
	}

	if refURL, err := url.Parse(redirectURL); err != nil || refURL.Host != request.Request().Host {
		u, _ := l.authManager.URL(ctx, "")
		redirectURL = u.String()
	}

	if redirectURL != "" {
		request.Session().Store("auth.redirect", redirectURL)
	}

	if l.loginTemplate != "" {
		return l.responder.Render(l.loginTemplate, nil)
	}

	return l.responder.RouteRedirect("auth.callback", nil)
}
package auth

type auth struct {
	Handler    handler
	Middleware middleware
}

func New() auth {
	authService := newService()

	return auth{
		Middleware: newMiddleware(authService),
		Handler:    newHandler(authService),
	}
}

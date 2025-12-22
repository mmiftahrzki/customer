package auth

type auth struct {
	Handler    handler
	Middleware middleware
}

func New(signingKey []byte) auth {
	service := newService(signingKey)
	handler := newHandler(service)

	return auth{
		Middleware: newMiddleware(service),
		Handler:    handler,
	}
}

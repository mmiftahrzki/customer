package auth

type auth struct {
	Handler    handler
	Middleware middleware
}

func New(signingKey []byte, adminEmail string) auth {
	service := newService(signingKey, adminEmail)
	handler := newHandler(service)

	return auth{
		Middleware: newMiddleware(service),
		Handler:    handler,
	}
}

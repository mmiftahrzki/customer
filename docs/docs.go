package docs

type docs struct {
	Handler handler
}

func New() docs {
	return docs{Handler: handler{}}
}

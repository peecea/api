package route

import (
	"peec/internal/route/api"
	"peec/internal/route/docs"
)

var RootRoutesGroup = []docs.RootDocumentation{
	{
		Group: "/api",
		Paths: api.Routes,
	},
}

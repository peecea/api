package route

import (
	"duval/internal/route/api"
	"duval/internal/route/docs"
)

var RootRoutesGroup = []docs.RootDocumentation{
	{
		Group: "/api",
		Paths: api.Routes,
	},
}

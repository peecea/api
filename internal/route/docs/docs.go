package docs

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RootDocumentation struct {
	Group string               `json:"group"`
	Paths []RouteDocumentation `json:"paths"`
}

type RouteDocumentation struct {
	HttpMethod   string          `json:"http_method"`
	RelativePath string          `json:"relative_path"`
	NeedToken    bool            `json:"need_token"`
	Handler      gin.HandlerFunc `json:"-"`
	DocRoot      string          `json:"-"`
}

func ParseDocumentation(documents []RootDocumentation) []RootDocumentation {
	var newDocuments []RootDocumentation

	for i := 0; i < len(documents); i++ {
		var newRoot RootDocumentation
		for j := 0; j < len(documents[i].Paths); j++ {
			newRoot.Paths = append(newRoot.Paths, documents[i].Paths[j])
		}
		newRoot.Group = documents[i].Group
		newDocuments = append(newDocuments, newRoot)
	}

	return newDocuments
}

func GenerateDocumentation(group *gin.RouterGroup, documents []RouteDocumentation) (err error) {
	for i := 0; i < len(documents); i++ {
		switch documents[i].HttpMethod {
		case http.MethodGet:
			group.GET(documents[i].RelativePath, documents[i].Handler)
			break
		case http.MethodPost:
			group.POST(documents[i].RelativePath, documents[i].Handler)
			break
		case http.MethodDelete:
			group.DELETE(documents[i].RelativePath, documents[i].Handler)
			break
		case http.MethodPut:
			group.PUT(documents[i].RelativePath, documents[i].Handler)
			break
		case http.MethodHead:
			group.Static(documents[i].RelativePath, documents[i].DocRoot)
			break
		}
	}

	return err
}

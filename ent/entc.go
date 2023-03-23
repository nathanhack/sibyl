//go:build ignore

package main

import (
	"fmt"
	"log"
	"net/http"

	"entgo.io/contrib/entoas"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"

	"ariga.io/ogent"
	"github.com/ogen-go/ogen"
)

//https://entgo.io/blog/2022/02/15/generate-rest-crud-with-ent-and-ogen/
//go generate ./...

func main() {
	spec := new(ogen.Spec)
	oas, err := entoas.NewExtension(
		entoas.Spec(spec),
		entoas.Mutations(func(graph *gen.Graph, spec *ogen.Spec) error {
			p := make(ogen.Paths)
			for k, v := range spec.Paths {
				p["/rest"+k] = v
			}
			spec.Paths = p

			spec.AddPathItem("/rest/search/{ticker}", ogen.NewPathItem().
				SetDescription("Searches for entities by ticker").
				SetGet(ogen.NewOperation().
					SetOperationID("searchTicker").
					SetSummary("Searches for entities by ticker.").
					AddTags("Search").
					AddResponse(fmt.Sprint(http.StatusOK), ogen.NewResponse().SetJSONContent(
						ogen.NewSchema().AddRequiredProperties(
							ogen.NewSchema().ToProperty("status"),
							ogen.NewSchema().AsArray().SetItems(
								ogen.NewSchema().AddRequiredProperties(
									ogen.String().ToProperty("ticker"),
									ogen.String().ToProperty("name"),
								),
							).ToProperty("results"),
							ogen.NewSchema().ToProperty("errors"),
						),
					)),
				).
				AddParameters(ogen.NewParameter().
					InPath().
					SetName("ticker").
					SetRequired(true).
					SetSchema(ogen.String()),
				),
			)

			spec.AddPathItem("/rest/entities/add/{ticker}", ogen.NewPathItem().
				SetDescription("Requests adding stock ticker").
				SetPost(ogen.NewOperation().
					SetOperationID("addTicker").
					SetSummary("Queue for adding entities by ticker.").
					AddTags("Add").
					AddResponse(fmt.Sprint(http.StatusOK), ogen.NewResponse().SetDescription("Ticker request received")),
				).
				AddParameters(ogen.NewParameter().
					InPath().
					SetName("ticker").
					SetRequired(true).
					SetSchema(ogen.String()),
				),
			)

			spec.AddPathItem("/rest/entities/add/{ticker}", ogen.NewPathItem().
				SetDescription("Requests adding stock ticker").
				SetPost(ogen.NewOperation().
					SetOperationID("addTicker").
					SetSummary("Queue for adding entities by ticker.").
					AddTags("Add").
					AddResponse(fmt.Sprint(http.StatusOK), ogen.NewResponse().SetDescription("Ticker request received")),
				).
				AddParameters(ogen.NewParameter().
					InPath().
					SetName("ticker").
					SetRequired(true).
					SetSchema(ogen.String()),
				),
			)

			return nil
		}),
	)

	ogent, err := ogent.NewExtension(spec)
	if err != nil {
		log.Fatalf("creating ogent extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ogent, oas))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

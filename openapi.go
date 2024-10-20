package main

import (
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"

	"github.com/LibraMusic/LibraCore/config"
	"github.com/LibraMusic/LibraCore/logging"
	"github.com/LibraMusic/LibraCore/types"
	"github.com/LibraMusic/LibraCore/utils"
)

func V1OpenAPI3Spec() openapi3.T {
	spec := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       config.Conf.Application.SourceName + " API",
			Description: "REST APIs used for interacting with the ToDo Service",
			Version:     utils.LibraVersion,
			License: &openapi3.License{
				Name: "MIT",
				URL:  "https://opensource.org/licenses/MIT",
			},
			Contact: &openapi3.Contact{
				URL: "https://github.com/LibraMusic/backend",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL: config.Conf.Application.PublicURL + ":" + strconv.Itoa(config.Conf.Application.Port) + "/api/v1",
			},
		},
	}

	spec.Components = &openapi3.Components{}

	trackSchema, err := openapi3gen.NewSchemaRefForValue(&types.Track{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for Track")
	}

	albumSchema, err := openapi3gen.NewSchemaRefForValue(&types.Album{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for Album")
	}

	videoSchema, err := openapi3gen.NewSchemaRefForValue(&types.Video{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for Video")
	}

	artistSchema, err := openapi3gen.NewSchemaRefForValue(&types.Artist{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for Artist")
	}

	playlistSchema, err := openapi3gen.NewSchemaRefForValue(&types.Playlist{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for Playlist")
	}

	userSchema, err := openapi3gen.NewSchemaRefForValue(&types.User{}, nil)
	if err != nil {
		logging.Error().Err(err).Msg("Failed to generate OpenAPI schema for User")
	}

	playableSchema := openapi3.NewOneOfSchema(
		trackSchema.Value,
		albumSchema.Value,
		videoSchema.Value,
		artistSchema.Value,
		playlistSchema.Value,
		userSchema.Value,
	).NewRef()

	spec.Components.Schemas = openapi3.Schemas{
		"Track":    trackSchema,
		"Album":    albumSchema,
		"Video":    videoSchema,
		"Artist":   artistSchema,
		"Playlist": playlistSchema,
		"User":     userSchema,
		"Playable": playableSchema,
	}

	spec.Components.Parameters = openapi3.ParametersMap{
		"SearchQueryParameter": &openapi3.ParameterRef{
			Value: openapi3.NewQueryParameter("q").
				WithDescription("Search query").
				WithRequired(true).
				WithSchema(openapi3.NewStringSchema()),
		},
	}

	spec.Components.RequestBodies = openapi3.RequestBodies{
		//
	}

	spec.Components.Responses = openapi3.ResponseBodies{
		"PlayableListResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returned when getting a list of playables").
				WithContent(openapi3.NewContentWithJSONSchema(openapi3.NewSchema().
					WithPropertyRef("playables", &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{openapi3.TypeArray},
							Items: &openapi3.SchemaRef{
								Ref: "#/components/schemas/Playable",
							},
						},
					}),
				)),
		},
	}

	spec.Paths = openapi3.NewPaths(
		openapi3.WithPath("/playables", &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "getAllPlayables",
				Summary:     "Gets all playables",
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(200, &openapi3.ResponseRef{
						Ref: "#/components/responses/PlayableListResponse",
					}),
				),
			},
		}),
		openapi3.WithPath("/search", &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "searchPlayables",
				Summary:     "Search for playables",
				Parameters: openapi3.Parameters{
					&openapi3.ParameterRef{
						Ref: "#/components/parameters/SearchQueryParameter",
					},
				},
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(200, &openapi3.ResponseRef{
						Ref: "#/components/responses/PlayableListResponse",
					}),
				),
			},
		}),
	)

	return spec
}

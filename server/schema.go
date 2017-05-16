package server

import (
	"errors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"log"
)

/* Schema in GraphQL spec language, used to guide the implementation below.

interface Node {
  id: ID!
}

type User : Node {
  id: ID!
  name: String
}

type Query {
  users: UserConnection
  node(id: ID!): Node
}
*/

// declare definitions first, initialize them in init()

var nodeDefinitions *relay.NodeDefinitions
var userType *graphql.Object

// Schema is the object representing our GraphQL schema
var Schema graphql.Schema

func init() {
	/**
	 * We get the node interface and field from the relay library.
	 *
	 * The first method is the way we resolve an ID to its object. The second is the
	 * way we resolve an object that implements node to its type.
	 */
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)

			// based on id and it stype, return the object
			switch resolvedID.Type {
			case "User":
				return GetUser(resolvedID.ID), nil
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			// based on the type of the value, return GraphQLObjectType
			switch p.Value.(type) {
			case *User:
				return userType
			default:
				return userType
			}
		},
	})

	// We define the user type
	userType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A webcig user.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("User", nil),
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the user.",
			},
		},
		Interfaces: []*graphql.Interface{nodeDefinitions.NodeInterface},
	})

	// We define a connection between a user and objects with a "users" field
	userConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "User",
		NodeType: userType,
	})

	// Now the root query type, the entry point into our schema
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: userConnectionDefinition.ConnectionType,
				Args: relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					log.Printf("Resolving users on root query.")
					args := relay.NewConnectionArguments(p.Args)
					users := []interface{}{}
					for _, user := range Users {
						users = append(users, user)
						log.Printf("%v\n", user)
					}
					return relay.ConnectionFromArray(users, args), nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	// Finally, we construct our schema
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		panic(err)
	}
}

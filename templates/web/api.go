package web

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "{{gomod.name}}/{{app.name}}/v1"
	"github.com/panyam/oneauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type {{app.classname}}Api struct {
	GrpcAddr       string
	LoggedInUserIdParam string
	DialOpts       []grpc.DialOption
	Context        context.Context
	AuthMiddleware *oneauth.AuthMiddleware
}

func (web *{{app.classname}}Api) SetupRoutes(apiRouter *mux.Router) {
	log.Println("Creating {{app.classname}}Api Router...")
	ctx := web.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Setup the API routes
	svcMux, err := web.createGatewayMux(ctx)
	if err != nil {
		log.Fatal(err)
	}
	v1 := apiRouter.PathPrefix("/v1")
	grpcgw := v1.PathPrefix("/{grpc_gateway:.*}").Subrouter()
	grpcgw.Handle("", http.StripPrefix("/api", svcMux))
}

// Creates the grpc runtime ServeMux that will ultimately forward all calls to /api/* to respective
// grpc endpoints defined in our protos
func (c *{{app.classname}}Api) createGatewayMux(ctx context.Context) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			//
			// Step 2 - Extend the context
			//
			metadata.AppendToOutgoingContext(ctx)

			loggedInUserId := c.AuthMiddleware.GetLoggedInUserId(request)

			md := metadata.Pairs()
			username, _, ok := request.BasicAuth()
			if ok { // TODO: Enable if turning on basic auth
				// TODO - Validate password if doing basic auth
				// send 401 if password is invalid
				md.Append("{{appname}}UserId", username)
			} else if loggedInUserId != "" {
				log.Println("LoggedInUserId: ", loggedInUserId)
				md.Append("{{appname}}UserId", loggedInUserId)
			}
			return md
		}))

	// Use the OpenTelemetry gRPC client interceptor for tracing
	conn, err := grpc.NewClient(c.GrpcAddr, c.DialOpts...)
	if err != nil {
		return nil, err
	}

{{ for service in services }}
	err = v1.Register{{service.name}}ServiceHandler(ctx, mux, conn)
	if err != nil {
		return nil, err
	}
{{ end }}
	return mux, nil
}

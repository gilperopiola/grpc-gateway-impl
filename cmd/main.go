package main

/* T0D0: Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests / Logging / Metrics / Tracing / Security / Error handling / Caching / Rate limiting / Postman collection */

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"strings"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

/* - Welcome~! - Here begins this simple implementation of the grpc-gateway framework. With gRPC, we design our service in a .proto file and then the server and client code is automatically generated. */

func main() {
	var (
		grpcServer  = initGRPCServer(v1Service.NewService())
		httpGateway = initHTTPGateway()
	)

	go runGRPCServer(grpcServer)
	runHTTPGateway(httpGateway)
}

/* So here we simulate a simple gRPC & HTTP backend API with 2 mock endpoints: Signup and Login.

 * With grpc-gateway we expose our gRPC service as a RESTful HTTP API, defining routes and verbs with annotations on the .proto files.
 * Then we just generate the gateway code and run it alongside the gRPC server. The gateway will translate HTTP requests to gRPC calls, handling input automatically.
 * We also use protovalidate to define input rules on the .proto files themselves for each request, which we enforce using an interceptor.

 * After the validation interceptor runs, requests go through one of the methods on pkg/v1/api.go. From there the Service layer is called, the business logic is executed and the request returns.
 */

const (
	gRPCPort = ":50051"
	httpPort = ":8080"

	errMsgListenGRPC = "Failed to listen gRPC: %v"
	errMsgServeGRPC  = "Failed to serve gRPC: %v"
	errMsgServeHTTP  = "Failed to serve HTTP: %v"
	errMsgGateway    = "Failed to start HTTP gateway: %v"
)

/* ----------------------------------- */
/*             - gRPC -                */
/* ----------------------------------- */

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
func initGRPCServer(service v1Service.ServiceLayer) *grpc.Server {
	var (
		interceptors = grpc.ChainUnaryInterceptor(
			v1.NewHTTPErrorHandlerInterceptor(),
			v1.NewValidationInterceptor(),
		)
		grpcServer = grpc.NewServer(interceptors)
	)

	usersPB.RegisterUsersServiceServer(grpcServer, &v1.API{Service: service})
	return grpcServer
}

// runGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server) {
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf(errMsgListenGRPC, err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(errMsgServeGRPC, err)
	}
}

/* ----------------------------------- */
/*             - HTTP -                */
/* ----------------------------------- */

func errorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	const fallback = `{"code": 13, "message": "failed to marshal error message"}` // return Internal when Marshal failed

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)
	pb := s.Proto()

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")

	contentType := marshaler.ContentType(pb)
	w.Header().Set("Content-Type", contentType)

	if s.Code() == codes.Unauthenticated {
		w.Header().Set("WWW-Authenticate", s.Message())
	}

	buf, merr := marshaler.Marshal(pb)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, mux, md)

	// RFC 7230 https://tools.ietf.org/html/rfc7230#section-4.1.2
	// Unless the request includes a TE header field indicating "trailers"
	// is acceptable, as described in Section 4.3, a server SHOULD NOT
	// generate trailer fields that it believes are necessary for the user
	// agent to receive.
	doForwardTrailers := requestAcceptsTrailers(r)

	if doForwardTrailers {
		handleForwardResponseTrailerHeader(w, mux, md)
		w.Header().Set("Transfer-Encoding", "chunked")
	}

	st := runtime.HTTPStatusFromCode(s.Code())
	if customStatus != nil {
		st = customStatus.HTTPStatus
	}

	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	if doForwardTrailers {
		handleForwardResponseTrailer(w, mux, md)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		if h, ok := defaultOutgoingTrailerMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		if h, ok := defaultOutgoingHeaderMatcher(k); ok {
			for _, v := range vs {
				w.Header().Add(h, v)
			}
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, mux *runtime.ServeMux, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		if h, ok := defaultOutgoingTrailerMatcher(k); ok {
			w.Header().Add("Trailer", textproto.CanonicalMIMEHeaderKey(h))
		}
	}
}

func defaultOutgoingHeaderMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

func defaultOutgoingTrailerMatcher(key string) (string, bool) {
	return fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, key), true
}

func requestAcceptsTrailers(req *http.Request) bool {
	te := req.Header.Get("TE")
	return strings.Contains(strings.ToLower(te), "trailers")
}

// initHTTPGateway initializes the HTTP gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
func initHTTPGateway() *runtime.ServeMux {
	mux := runtime.NewServeMux(runtime.WithErrorHandler(errorHandler))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, gRPCPort, opts); err != nil {
		log.Fatalf(errMsgGateway, err)
	}

	return mux
}

// runHTTPGateway runs the HTTP gateway on a given port.
// It listens for incoming HTTP requests and serves them.
func runHTTPGateway(httpGateway *runtime.ServeMux) {
	log.Println("Running HTTP!")

	if err := http.ListenAndServe(httpPort, httpGateway); err != nil {
		log.Fatalf(errMsgServeHTTP, err)
	}
}

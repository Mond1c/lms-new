package obs

import (
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func WrapHTTP(h http.Handler, serviceName string) http.Handler {
	return otelhttp.NewHandler(h, serviceName+"-http")
}

func ConnectInterceptor() (connect.Interceptor, error) {
	return otelconnect.NewInterceptor(otelconnect.WithTrustRemote())
}

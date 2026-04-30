# Figaro

[![Build Status](https://github.com/alexdebril/figaro/actions/workflows/go.yml/badge.svg)](https://github.com/alexdebril/figaro/actions/workflows/go.yml/)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexdebril/figaro)](https://goreportcard.com/report/github.com/alexdebril/figaro)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=alexdebril_figaro&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=alexdebril_figaro)

Small library to make easy Go 

*Easy come, easy go, will you let me go?*

The library takes its name from a famous song wrote by a famous artist as the goal of those packages is to make "easy go".

## Packages

### http/response

This package provides helpers designed to help in writing HTTP response. The principle is easy: create a new Response fed with a set of options and inject the `http.responseWriter` into it to write the response.

```go
package your_package

import (
	"net/http"
	"time"

	"github.com/alexdebril/figaro/api/response"
)

type JsonMessage struct {
    Date    time.Time `json:"date"`
    Message string    `json:"message"`
}

type Handler struct {}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    message := &JsonMessage{
        Date:    time.Now(),
        Message: "this is the new message",
    }
    resp := response.NewJsonResponse(message)
    resp.Write(w)
}
```

### Telemetry

Figaro provides the `otel` package to help in writing telemetry. 

It supports metrics and traces. Logs are excluded from the scope as it's not stable yet.

#### HTTP Handler

use the following code to wrap your HTTP handler:

```go
// build a context with custom attributes to define the resource's name and version
ctx := newContext()
// build the API's Mux in a factory function ...
mux := newAPIMux(ctx, l)
// wrap it with OTEL's tracing Mux ...
handler, shutdown := otel.NewHTTPHandler(ctx, l, mux, otel.WithTracing(), otel.WithMetrics())
defer shutdown(ctx)
// and finally run the server.
s := api.NewServer(l, handler, f.address)
if err := s.Listen(ctx); err != nil {
    l.Error("server returned an error", slog.Any("err", err))
    os.Exit(1)
}
```

See the [example](https://github.com/alexdebril/figaro/blob/main/examples/simple-api-otel/main.go) for more details.

#### Instrumenting your code

Everything's based on the [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go). Visit the [Instrumentation](https://opentelemetry.io/docs/languages/go/instrumentation/) section to see how to instrument your code.

#### Testing telemetry

It's easy to test if metrics and traces are being sent to the right place.

```shell
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```
Run your app, and then go to http://localhost:16686 to see the traces.

#### Auto Exporters

Figaro uses auto exporters to send metrics and traces to the right place. Its configuration is done by environment variables, as documented in the [OpenTelemetry Exporters Section](https://opentelemetry.io/docs/languages/go/exporters/#automatic-exporter-configuration-with-environment-variables)

For instance:

```shell
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
```
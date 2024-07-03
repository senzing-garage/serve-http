package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docktermj/cloudshell/xtermservice"
	"github.com/flowchartsman/swaggerui"
	"github.com/pkg/browser"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestapi"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	APIUrlRoutePrefix         string // FIXME: Only works with "api"
	AvoidServing              bool
	EnableAll                 bool
	EnableSenzingRestAPI      bool
	EnableSwaggerUI           bool
	EnableXterm               bool
	GrpcDialOptions           []grpc.DialOption
	GrpcTarget                string
	LogLevelName              string
	ObserverOrigin            string
	Observers                 []observer.Observer
	OpenAPISpecificationRest  []byte
	ReadHeaderTimeout         time.Duration
	SenzingSettings           string
	SenzingInstanceName       string
	SenzingVerboseLogging     int64
	ServerAddress             string
	ServerOptions             []senzingrestapi.ServerOption
	ServerPort                int
	SwaggerURLRoutePrefix     string // FIXME: Only works with "swagger"
	TtyOnly                   bool
	XtermAllowedHostnames     []string
	XtermArguments            []string
	XtermCommand              string
	XtermConnectionErrorLimit int
	XtermKeepalivePingTimeout int
	XtermMaxBufferSizeBytes   int
	XtermURLRoutePrefix       string // FIXME: Only works with "xterm"
}

type TemplateVariables struct {
	BasicHTTPServer
	APIServerStatus string
	APIServerURL    string
	HTMLTitle       string
	RequestHost     string
	SwaggerStatus   string
	SwaggerURL      string
	XtermStatus     string
	XtermURL        string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method simply prints the 'Something' value in the type-struct.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	var err error
	var userMessage = ""

	rootMux := http.NewServeMux()

	// Enable Senzing HTTP REST API.

	if httpServer.EnableAll || httpServer.EnableSenzingRestAPI {
		senzingAPIMux := httpServer.getSenzingAPIMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.APIUrlRoutePrefix), http.StripPrefix("/api", senzingAPIMux))
		userMessage = fmt.Sprintf("%sServing Senzing REST API at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.APIUrlRoutePrefix)
	}

	// Enable SwaggerUI.

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUIMux := httpServer.getSwaggerUIMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.SwaggerURLRoutePrefix), http.StripPrefix("/swagger", swaggerUIMux))
		userMessage = fmt.Sprintf("%sServing SwaggerUI at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.SwaggerURLRoutePrefix)
	}

	// Enable Xterm.

	if httpServer.EnableAll || httpServer.EnableXterm {
		err := os.Setenv("SENZING_ENGINE_CONFIGURATION_JSON", httpServer.SenzingSettings)
		if err != nil {
			panic(err)
		}
		xtermMux := httpServer.getXtermMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.XtermURLRoutePrefix), http.StripPrefix("/xterm", xtermMux))
		userMessage = fmt.Sprintf("%sServing XTerm at            http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.XtermURLRoutePrefix)
	}

	// Add route to template pages.

	rootMux.HandleFunc("/site/", httpServer.siteFunc)
	userMessage = fmt.Sprintf("%sServing Console at          http://localhost:%d\n", userMessage, httpServer.ServerPort)

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}
	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessage = fmt.Sprintf("%sStarting server on interface:port '%s'...\n", userMessage, listenOnAddress)
	fmt.Println(userMessage)
	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}

	// Start a web browser.  Unless disabled.

	if !httpServer.TtyOnly {
		_ = browser.OpenURL(fmt.Sprintf("http://localhost:%d", httpServer.ServerPort))
	}

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()
	}
	return err
}

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) getServerStatus(up bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}
	if up {
		result = "green"
	}
	return result
}

func (httpServer *BasicHTTPServer) getServerURL(up bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}
	if up {
		result = url
	}
	return result
}

func (httpServer *BasicHTTPServer) openAPIFunc(ctx context.Context, openAPISpecification []byte) http.HandlerFunc {
	_ = ctx
	_ = openAPISpecification
	return func(w http.ResponseWriter, r *http.Request) {
		var bytesBuffer bytes.Buffer
		bufioWriter := bufio.NewWriter(&bytesBuffer)
		openAPISpecificationTemplate, err := template.New("OpenApiTemplate").Parse(string(httpServer.OpenAPISpecificationRest))
		if err != nil {
			panic(err)
		}
		templateVariables := TemplateVariables{
			RequestHost: r.Host,
		}
		err = openAPISpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
func (httpServer *BasicHTTPServer) populateStaticTemplate(responseWriter http.ResponseWriter, request *http.Request, filepath string, templateVariables TemplateVariables) {
	_ = request
	templateBytes, err := static.ReadFile(filepath)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	err = templateParsed.Execute(responseWriter, templateVariables)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *BasicHTTPServer) getSenzingAPIMux(ctx context.Context) *senzingrestapi.Server {
	_ = ctx
	service := &senzingrestservice.BasicSenzingRestService{
		GrpcDialOptions:          httpServer.GrpcDialOptions,
		GrpcTarget:               httpServer.GrpcTarget,
		LogLevelName:             httpServer.LogLevelName,
		ObserverOrigin:           httpServer.ObserverOrigin,
		Observers:                httpServer.Observers,
		Settings:                 httpServer.SenzingSettings,
		SenzingInstanceName:      httpServer.SenzingInstanceName,
		SenzingVerboseLogging:    httpServer.SenzingVerboseLogging,
		URLRoutePrefix:           httpServer.APIUrlRoutePrefix,
		OpenAPISpecificationSpec: httpServer.OpenAPISpecificationRest,
	}
	srv, err := senzingrestapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

func (httpServer *BasicHTTPServer) getSwaggerUIMux(ctx context.Context) *http.ServeMux {
	swaggerMux := swaggerui.Handler([]byte{}) // OpenAPI specification handled by openApiFunc()
	swaggerFunc := swaggerMux.ServeHTTP
	submux := http.NewServeMux()
	submux.HandleFunc("/", swaggerFunc)
	submux.HandleFunc("/swagger_spec", httpServer.openAPIFunc(ctx, httpServer.OpenAPISpecificationRest))
	return submux
}

func (httpServer *BasicHTTPServer) getXtermMux(ctx context.Context) *http.ServeMux {
	xtermService := &xtermservice.XtermServiceImpl{
		AllowedHostnames:     httpServer.XtermAllowedHostnames,
		Arguments:            httpServer.XtermArguments,
		Command:              httpServer.XtermCommand,
		ConnectionErrorLimit: httpServer.XtermConnectionErrorLimit,
		KeepalivePingTimeout: httpServer.XtermKeepalivePingTimeout,
		MaxBufferSizeBytes:   httpServer.XtermMaxBufferSizeBytes,
		UrlRoutePrefix:       httpServer.XtermURLRoutePrefix,
	}
	return xtermService.Handler(ctx)
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *BasicHTTPServer) siteFunc(w http.ResponseWriter, r *http.Request) {
	templateVariables := TemplateVariables{
		BasicHTTPServer: *httpServer,
		HTMLTitle:       "Senzing Tools",
		APIServerURL:    httpServer.getServerURL(httpServer.EnableSenzingRestAPI, fmt.Sprintf("http://%s/api", r.Host)),
		APIServerStatus: httpServer.getServerStatus(httpServer.EnableSenzingRestAPI),
		SwaggerURL:      httpServer.getServerURL(httpServer.EnableSwaggerUI, fmt.Sprintf("http://%s/swagger", r.Host)),
		SwaggerStatus:   httpServer.getServerStatus(httpServer.EnableSwaggerUI),
		XtermURL:        httpServer.getServerURL(httpServer.EnableXterm, fmt.Sprintf("http://%s/xterm", r.Host)),
		XtermStatus:     httpServer.getServerStatus(httpServer.EnableXterm),
	}
	w.Header().Set("Content-Type", "text/html")
	filePath := fmt.Sprintf("static/templates%s", r.RequestURI)
	httpServer.populateStaticTemplate(w, r, filePath, templateVariables)
}

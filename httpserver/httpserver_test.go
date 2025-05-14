package httpserver_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-rest-api-service/senzingrestservice"
	"github.com/senzing-garage/serve-http/httpserver"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_Serve(test *testing.T) {
	ctx := test.Context()
	httpServer := getTestObject(ctx, test)
	err := httpServer.Serve(ctx)
	require.NoError(test, err)
}

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

// func TestBasicHTTPServer_getServerStatus(test *testing.T) {
// 	_ = test
// 	ctx := context.TODO()
// 	httpServer := getTestObject(ctx, test)
// 	actual := httpServer.getServerStatus(true)
// 	assert.Equal(test, "green", actual)
// }

// func TestBasicHTTPServer_getServerURL(test *testing.T) {
// 	_ = test
// 	ctx := context.TODO()
// 	expected := "http://expected"
// 	httpServer := getTestObject(ctx, test)
// 	actual := httpServer.getServerURL(true, expected)
// 	assert.Equal(test, expected, actual)
// }

// func TestBasicHTTPServer_openAPIFunc(test *testing.T) {
// 	_ = test
// 	ctx := context.TODO()
// 	httpServer := getTestObject(ctx, test)
// 	openAPIFunction := httpServer.openAPIFunc(ctx, httpServer.OpenAPISpecificationRest)
// 	request := httptest.NewRequest(http.MethodGet, "/", nil)
// 	response := httptest.NewRecorder()
// 	openAPIFunction(response, request)
// }

// func TestBasicHTTPServer_populateStaticTemplate(test *testing.T) {
// 	_ = test
// 	ctx := context.TODO()
// 	request := httptest.NewRequest(http.MethodGet, "/", nil)
// 	response := httptest.NewRecorder()
// 	httpServer := getTestObject(ctx, test)
// 	httpServer.populateStaticTemplate(response, request, "/", httpserver.TemplateVariables{})
// }

// func TestBasicHTTPServer_siteFunc(test *testing.T) {
// 	_ = test
// 	ctx := context.TODO()
// 	request := httptest.NewRequest(http.MethodGet, "/", nil)
// 	response := httptest.NewRecorder()
// 	httpServer := getTestObject(ctx, test)
// 	httpServer.siteFunc(response, request)
// }

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, t *testing.T) *httpserver.BasicHTTPServer {
	t.Helper()

	_ = ctx

	observer1 := &observer.NullObserver{
		ID: "Observer 1",
	}

	logLevelName := "INFO"
	osenvLogLevel := os.Getenv("SENZING_LOG_LEVEL")

	if len(osenvLogLevel) > 0 {
		logLevelName = osenvLogLevel
	}

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &httpserver.BasicHTTPServer{
		APIUrlRoutePrefix:        "api",
		AvoidServing:             true,
		EnableAll:                true,
		LogLevelName:             logLevelName,
		ObserverOrigin:           "Test Observer origin",
		Observers:                []observer.Observer{observer1},
		OpenAPISpecificationRest: senzingrestservice.OpenAPISpecificationJSON,
		ReadHeaderTimeout:        10 * time.Second,
		SenzingInstanceName:      "Test HTTP Server",
		SenzingSettings:          senzingSettings,
		SwaggerURLRoutePrefix:    "swagger",
		TtyOnly:                  true,
		XtermURLRoutePrefix:      "xterm",
	}

	return result
}

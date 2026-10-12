package httpserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/senzing-garage/go-helpers/settings"
	"github.com/senzing-garage/serve-http/httpserver"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicHTTPServer_Serve(test *testing.T) {
	test.Parallel()
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

	senzingSettings, err := settings.BuildSimpleSettingsUsingEnvVars()
	require.NoError(t, err)

	result := &httpserver.BasicHTTPServer{
		AvoidServing:        true,
		EnableAll:           true,
		ReadHeaderTimeout:   10 * time.Second,
		SenzingSettings:     senzingSettings,
		TtyOnly:             true,
		XtermURLRoutePrefix: "xterm",
	}

	return result
}

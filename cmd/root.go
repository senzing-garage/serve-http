/*
 */
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-grpcing/grpcurl"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-rest-api-service/senzingrestservice"
	"github.com/senzing/senzing-tools/cmdhelper"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/help"
	"github.com/senzing/senzing-tools/option"
	"github.com/senzing/serve-http/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	Short string = "HTTP server supporting various services"
	Use   string = "serve-http"
	Long  string = `
An HTTP server supporting the following services:
	- Senzing API server
	- Swagger UI
	- Xterm
	`
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextBools = []cmdhelper.ContextBool{
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableSwaggerUi, false),
		Envar:   envar.EnableSwaggerUi,
		Help:    help.EnableSwaggerUi,
		Option:  option.EnableSwaggerUi,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableAll, false),
		Envar:   envar.EnableAll,
		Help:    help.EnableAll,
		Option:  option.EnableAll,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableSenzingRestApi, false),
		Envar:   envar.EnableSenzingRestApi,
		Help:    help.EnableSenzingRestApi,
		Option:  option.EnableSenzingRestApi,
	},
	{
		Default: cmdhelper.OsLookupEnvBool(envar.EnableXterm, false),
		Envar:   envar.EnableXterm,
		Help:    help.EnableXterm,
		Option:  option.EnableXterm,
	},
}

var ContextInts = []cmdhelper.ContextInt{
	{
		Default: cmdhelper.OsLookupEnvInt(envar.Configuration, 0),
		Envar:   envar.EngineLogLevel,
		Help:    help.EngineLogLevel,
		Option:  option.EngineLogLevel,
	},
	{
		Default: cmdhelper.OsLookupEnvInt(envar.HttpPort, 8261),
		Envar:   envar.HttpPort,
		Help:    help.HttpPort,
		Option:  option.HttpPort,
	},
	{
		Default: cmdhelper.OsLookupEnvInt(envar.XtermConnectionErrorLimit, 10),
		Envar:   envar.XtermConnectionErrorLimit,
		Help:    help.XtermConnectionErrorLimit,
		Option:  option.XtermConnectionErrorLimit,
	},
	{
		Default: cmdhelper.OsLookupEnvInt(envar.XtermKeepalivePingTimeout, 20),
		Envar:   envar.XtermKeepalivePingTimeout,
		Help:    help.XtermKeepalivePingTimeout,
		Option:  option.XtermKeepalivePingTimeout,
	},
	{
		Default: cmdhelper.OsLookupEnvInt(envar.XtermMaxBufferSizeBytes, 512),
		Envar:   envar.XtermMaxBufferSizeBytes,
		Help:    help.XtermMaxBufferSizeBytes,
		Option:  option.XtermMaxBufferSizeBytes,
	},
}

var ContextStrings = []cmdhelper.ContextString{
	{
		Default: cmdhelper.OsLookupEnvString(envar.Configuration, ""),
		Envar:   envar.Configuration,
		Help:    help.Configuration,
		Option:  option.Configuration,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.DatabaseUrl, ""),
		Envar:   envar.DatabaseUrl,
		Help:    help.DatabaseUrl,
		Option:  option.DatabaseUrl,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.EngineConfigurationJson, ""),
		Envar:   envar.EngineConfigurationJson,
		Help:    help.EngineConfigurationJson,
		Option:  option.EngineConfigurationJson,
	},
	{
		Default: fmt.Sprintf("serve-http-%d", time.Now().Unix()),
		Envar:   envar.EngineModuleName,
		Help:    help.EngineModuleName,
		Option:  option.EngineModuleName,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.GrpcUrl, ""),
		Envar:   envar.GrpcUrl,
		Help:    help.GrpcUrl,
		Option:  option.GrpcUrl,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.LogLevel, "INFO"),
		Envar:   envar.LogLevel,
		Help:    help.LogLevel,
		Option:  option.LogLevel,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverOrigin, "serve-http"),
		Envar:   envar.ObserverOrigin,
		Help:    help.ObserverOrigin,
		Option:  option.ObserverOrigin,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverUrl, ""),
		Envar:   envar.ObserverUrl,
		Help:    help.ObserverUrl,
		Option:  option.ObserverUrl,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ServerAddress, "0.0.0.0"),
		Envar:   envar.ServerAddress,
		Help:    help.ServerAddress,
		Option:  option.ServerAddress,
	},
	{
		Default: cmdhelper.OsLookupEnvString(envar.ObserverUrl, "/bin/bash"),
		Envar:   envar.XtermCommand,
		Help:    help.XtermCommand,
		Option:  option.XtermCommand,
	},
}

var ContextStringSlices = []cmdhelper.ContextStringSlice{
	{
		Default: getDefaultAllowedHostnames(),
		Envar:   envar.XtermAllowedHostnames,
		Help:    help.XtermAllowedHostnames,
		Option:  option.XtermAllowedHostnames,
	},
	{
		Default: []string{},
		Envar:   envar.XtermArguments,
		Help:    help.XtermArguments,
		Option:  option.XtermArguments,
	},
}

var ContextVariables = &cmdhelper.ContextVariables{
	Bools:   append(ContextBools, ContextBoolsForOsArch...),
	Ints:    append(ContextInts, ContextIntsForForOsArch...),
	Strings: append(ContextStrings, ContextStringsForOsArch...),
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, *ContextVariables)
}

// --- Networking -------------------------------------------------------------

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func getDefaultAllowedHostnames() []string {
	result := []string{"localhost"}
	outboundIpAddress := getOutboundIP().String()
	if len(outboundIpAddress) > 0 {
		result = append(result, outboundIpAddress)
	}
	return result
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, *ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	// Build senzingEngineConfigurationJson.

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		options := map[string]string{
			"configPath":          viper.GetString(option.ConfigPath),
			"databaseUrl":         viper.GetString(option.DatabaseUrl),
			"licenseStringBase64": viper.GetString(option.LicenseStringBase64),
			"resourcePath":        viper.GetString(option.ResourcePath),
			"senzingDirectory":    viper.GetString(option.SenzingDirectory),
			"supportPath":         viper.GetString(option.SupportPath),
		}
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJsonUsingMap(options)
		if err != nil {
			return err
		}
	}
	err = g2engineconfigurationjson.VerifySenzingEngineConfigurationJson(ctx, senzingEngineConfigurationJson)
	if err != nil {
		return err
	}

	// Determine if gRPC is being used.

	grpcUrl := viper.GetString(option.GrpcUrl)
	grpcTarget := ""
	grpcDialOptions := []grpc.DialOption{}
	if len(grpcUrl) > 0 {
		grpcTarget, grpcDialOptions, err = grpcurl.Parse(ctx, grpcUrl)
		if err != nil {
			return err
		}
	}

	// Build observers.
	//  viper.GetString(option.ObserverUrl),

	observers := []observer.Observer{}

	// Create object and Serve.

	httpServer := &httpserver.HttpServerImpl{
		ApiUrlRoutePrefix:              "api",
		EnableAll:                      viper.GetBool(option.EnableAll),
		EnableSenzingRestAPI:           viper.GetBool(option.EnableSenzingRestApi),
		EnableSwaggerUI:                viper.GetBool(option.EnableSwaggerUi),
		EnableXterm:                    viper.GetBool(option.EnableXterm),
		GrpcDialOptions:                grpcDialOptions,
		GrpcTarget:                     grpcTarget,
		LogLevelName:                   viper.GetString(option.LogLevel),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin),
		Observers:                      observers,
		OpenApiSpecificationRest:       senzingrestservice.OpenApiSpecificationJson,
		ReadHeaderTimeout:              60 * time.Second,
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
		ServerAddress:                  viper.GetString(option.ServerAddress),
		ServerPort:                     viper.GetInt(option.HttpPort),
		SwaggerUrlRoutePrefix:          "swagger",
		XtermAllowedHostnames:          viper.GetStringSlice(option.XtermAllowedHostnames),
		XtermArguments:                 viper.GetStringSlice(option.XtermArguments),
		XtermCommand:                   viper.GetString(option.XtermCommand),
		XtermConnectionErrorLimit:      viper.GetInt(option.XtermConnectionErrorLimit),
		XtermKeepalivePingTimeout:      viper.GetInt(option.XtermKeepalivePingTimeout),
		XtermMaxBufferSizeBytes:        viper.GetInt(option.XtermMaxBufferSizeBytes),
		XtermUrlRoutePrefix:            "xterm",
	}
	return httpServer.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}

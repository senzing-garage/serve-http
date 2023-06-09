/*
 */
package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-grpcing/grpcurl"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-rest-api-service/senzingrestservice"
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/help"
	"github.com/senzing/senzing-tools/helper"
	"github.com/senzing/senzing-tools/option"
	"github.com/senzing/serve-http/httpserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	defaultConfiguration             string = ""
	defaultDatabaseUrl               string = ""
	defaultEnableAll                 bool   = false
	defaultEnableSenzingRestApi      bool   = false
	defaultEnableSwaggerUI           bool   = false
	defaultEnableXterm               bool   = false
	defaultEngineConfigurationJson   string = ""
	defaultEngineLogLevel            int    = 0
	defaultGrpcUrl                          = ""
	defaultHttpPort                  int    = 8261
	defaultLogLevel                  string = "INFO"
	defaultObserverOrigin            string = "serve-http"
	defaultObserverUrl               string = ""
	defaultServerAddress             string = "0.0.0.0"
	defaultXtermCommand              string = "/bin/bash"
	defaultXtermConnectionErrorLimit int    = 10
	defaultXtermKeepalivePingTimeout int    = 20
	defaultXtermMaxBufferSizeBytes   int    = 512
	// envarEnableAll                   string = "SENZING_TOOLS_ENABLE_ALL"
	// envarEnableSenzingRestApi        string = "SENZING_TOOLS_ENABLE_SENZING_REST_API"
	// envarEnableXterm                 string = "SENZING_TOOLS_ENABLE_XTERM"
	// envarServerAddress               string = "SENZING_TOOLS_SERVER_ADDRESS"
	// envarXtermAllowedHostnames       string = "SENZING_TOOLS_XTERM_ALLOWED_HOSTNAMES"
	// envarXtermArguments              string = "SENZING_TOOLS_XTERM_ARGUMENTS"
	// envarXtermCommand                string = "SENZING_TOOLS_XTERM_COMMAND"
	// envarXtermConnectionErrorLimit   string = "SENZING_TOOLS_XTERM_CONNECTION_ERROR_LIMIT"
	// envarXtermKeepalivePingTimeout   string = "SENZING_TOOLS_XTERM_KEEPALIVE_PING_TIMEOUT"
	// envarXtermMaxBufferSizeBytes     string = "SENZING_TOOLS_XTERM_MAX_BUFFER_SIZE_BYTES"
	// optionEnableAll                  string = "enable-all"
	// optionEnableSenzingRestApi       string = "enable-senzing-rest-api"
	// optionEnableXterm                string = "enable-xterm"
	// optionServerAddress              string = "server-address"
	// optionXtermAllowedHostnames      string = "xterm-allowed-hostnames"
	// optionXtermArguments             string = "xterm-arguments"
	// optionXtermCommand               string = "xterm-command"
	// optionXtermConnectionErrorLimit  string = "xterm-connection-error-limit"
	// optionXtermKeepalivePingTimeout  string = "xterm-keepalive-ping-timeout"
	// optionXtermMaxBufferSizeBytes    string = "xterm-max-buffer-size-bytes"
	Short string = "serve-http short description"
	Use   string = "serve-http"
	Long  string = `
serve-http long description.
	`
)

var (
	defaultEngineModuleName      string   = fmt.Sprintf("serve-http-%d", time.Now().Unix())
	defaultXtermAllowedHostnames []string = getDefaultAllowedHostnames()
	defaultXtermArguments        []string
)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	RootCmd.Flags().Bool(option.EnableSwaggerUi, defaultEnableSwaggerUI, fmt.Sprintf(help.EnableSwaggerUi, envar.EnableSwaggerUi))
	RootCmd.Flags().Bool(option.EnableAll, defaultEnableAll, fmt.Sprintf(help.EnableAll, envar.EnableAll))
	RootCmd.Flags().Bool(option.EnableSenzingRestApi, defaultEnableSenzingRestApi, fmt.Sprintf(help.EnableSenzingRestApi, envar.EnableSenzingRestApi))
	RootCmd.Flags().Bool(option.EnableXterm, defaultEnableXterm, fmt.Sprintf(help.EnableXterm, envar.EnableXterm))
	RootCmd.Flags().Int(option.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf(help.EngineLogLevel, envar.EngineLogLevel))
	RootCmd.Flags().Int(option.HttpPort, defaultHttpPort, fmt.Sprintf(help.HttpPort, envar.HttpPort))
	RootCmd.Flags().Int(option.XtermConnectionErrorLimit, defaultXtermConnectionErrorLimit, fmt.Sprintf(help.XtermConnectionErrorLimit, envar.XtermConnectionErrorLimit))
	RootCmd.Flags().Int(option.XtermKeepalivePingTimeout, defaultXtermKeepalivePingTimeout, fmt.Sprintf(help.XtermKeepalivePingTimeout, envar.XtermKeepalivePingTimeout))
	RootCmd.Flags().Int(option.XtermMaxBufferSizeBytes, defaultXtermMaxBufferSizeBytes, fmt.Sprintf(help.XtermMaxBufferSizeBytes, envar.XtermMaxBufferSizeBytes))
	RootCmd.Flags().String(option.Configuration, defaultConfiguration, fmt.Sprintf(help.Configuration, envar.Configuration))
	RootCmd.Flags().String(option.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf(help.DatabaseUrl, envar.DatabaseUrl))
	RootCmd.Flags().String(option.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf(help.EngineConfigurationJson, envar.EngineConfigurationJson))
	RootCmd.Flags().String(option.EngineModuleName, defaultEngineModuleName, fmt.Sprintf(help.EngineModuleName, envar.EngineModuleName))
	RootCmd.Flags().String(option.GrpcUrl, defaultGrpcUrl, fmt.Sprintf(help.GrpcUrl, envar.GrpcUrl))
	RootCmd.Flags().String(option.LogLevel, defaultLogLevel, fmt.Sprintf(help.LogLevel, envar.LogLevel))
	RootCmd.Flags().String(option.ObserverOrigin, defaultObserverOrigin, fmt.Sprintf(help.ObserverOrigin, envar.ObserverOrigin))
	RootCmd.Flags().String(option.ObserverUrl, defaultObserverUrl, fmt.Sprintf(help.ObserverUrl, envar.ObserverUrl))
	RootCmd.Flags().String(option.ServerAddress, defaultServerAddress, fmt.Sprintf(help.ServerAddress, envar.ServerAddress))
	RootCmd.Flags().String(option.XtermCommand, defaultXtermCommand, fmt.Sprintf(help.XtermCommand, envar.XtermCommand))
	RootCmd.Flags().StringSlice(option.XtermAllowedHostnames, defaultXtermAllowedHostnames, fmt.Sprintf(help.XtermAllowedHostnames, envar.XtermAllowedHostnames))
	RootCmd.Flags().StringSlice(option.XtermArguments, defaultXtermArguments, fmt.Sprintf(help.XtermArguments, envar.XtermArguments))
}

// If a configuration file is present, load it.
func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := ""
	configFlag := cobraCommand.Flags().Lookup(option.Configuration)
	if configFlag != nil {
		configuration = configFlag.Value.String()
	}
	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("serve-http")
		viper.SetConfigType("yaml")

		// Define search path order.

		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
	}

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Applying configuration file:", viper.ConfigFileUsed())
	}
}

// Configure Viper with user-specified options.
func loadOptions(cobraCommand *cobra.Command) {
	var err error = nil
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Bools

	boolOptions := map[string]bool{
		option.EnableAll:            defaultEnableAll,
		option.EnableSenzingRestApi: defaultEnableSenzingRestApi,
		option.EnableSwaggerUi:      defaultEnableSwaggerUI,
		option.EnableXterm:          defaultEnableXterm,
	}
	for optionKey, optionValue := range boolOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// Ints

	intOptions := map[string]int{
		option.EngineLogLevel:            defaultEngineLogLevel,
		option.HttpPort:                  defaultHttpPort,
		option.XtermConnectionErrorLimit: defaultXtermConnectionErrorLimit,
		option.XtermKeepalivePingTimeout: defaultXtermKeepalivePingTimeout,
		option.XtermMaxBufferSizeBytes:   defaultXtermMaxBufferSizeBytes,
	}
	for optionKey, optionValue := range intOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// Strings

	stringOptions := map[string]string{
		option.Configuration:           defaultConfiguration,
		option.DatabaseUrl:             defaultDatabaseUrl,
		option.EngineConfigurationJson: defaultEngineConfigurationJson,
		option.EngineModuleName:        defaultEngineModuleName,
		option.GrpcUrl:                 defaultGrpcUrl,
		option.LogLevel:                defaultLogLevel,
		option.ObserverOrigin:          defaultObserverOrigin,
		option.ObserverUrl:             defaultObserverUrl,
		option.ServerAddress:           defaultServerAddress,
		option.XtermCommand:            defaultXtermCommand,
	}
	for optionKey, optionValue := range stringOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// StringSlice

	stringSliceOptions := map[string][]string{
		option.XtermAllowedHostnames: defaultXtermAllowedHostnames,
		option.XtermArguments:        defaultXtermArguments,
	}
	for optionKey, optionValue := range stringSliceOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

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
	loadConfigurationFile(cobraCommand)
	loadOptions(cobraCommand)
	cobraCommand.SetVersionTemplate(constant.VersionTemplate)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.TODO()

	// Build senzingEngineConfigurationJson.

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(option.DatabaseUrl))
		if err != nil {
			return err
		}
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
	err = httpServer.Serve(ctx)
	return err
}

// Used in construction of cobra.Command
func Version() string {
	return helper.MakeVersion(githubVersion, githubIteration)
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

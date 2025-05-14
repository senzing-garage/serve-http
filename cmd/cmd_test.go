package cmd_test

import (
	"os"
	"testing"

	"github.com/senzing-garage/serve-http/cmd"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "--avoid-serving", "--tty-only"}

	cmd.Execute()
}

// func Test_Execute_completion(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "completion"}

// 	cmd.Execute()
// }

// func Test_Execute_docs(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "docs"}

// 	cmd.Execute()
// }

// func Test_Execute_help(test *testing.T) {
// 	test.Parallel()

// 	os.Args = []string{"command-name", "--help"}

// 	cmd.Execute()
// }

// func Test_PreRun(test *testing.T) {
// 	test.Parallel()

// 	args := []string{"command-name", "--help"}
// 	cmd.PreRun(cmd.RootCmd, args)
// }

// func Test_RunE(test *testing.T) {
// 	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")

// 	err := cmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_RootCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.RootCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.RootCmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_completionCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.CompletionCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.CompletionCmd.RunE(cmd.CompletionCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_docsCmd(test *testing.T) {
// 	test.Parallel()

// 	err := cmd.DocsCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.DocsCmd.RunE(cmd.DocsCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_RunE_badGrpcURL(test *testing.T) {
// 	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")
// 	test.Setenv("SENZING_TOOLS_GRPC_URL", "grpc://bad")

// 	err := cmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

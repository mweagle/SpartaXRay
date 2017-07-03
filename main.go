package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	spartaCGO "github.com/mweagle/Sparta/cgo"
	gocf "github.com/mweagle/go-cloudformation"
)

// Standard AWS λ function
func helloWorld(event *json.RawMessage,
	context *sparta.LambdaContext,
	w http.ResponseWriter,
	logger *logrus.Logger) {

	fmt.Fprint(w, "Hello World 🌍")
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	lambdaOptions := &sparta.LambdaFunctionOptions{
		MemorySize: 256,
		TracingConfig: &gocf.LambdaFunctionTracingConfig{
			Mode: gocf.String("Active"),
		},
		SpartaOptions: &sparta.SpartaOptions{
			Name: "SpartaXRay",
		},
	}
	lambdaOptions.Tags = map[string]string{"special": "tag"}
	lambdaFn := sparta.NewLambda(sparta.IAMRoleDefinition{},
		helloWorld,
		lambdaOptions)

	// Sanitize the name so that it doesn't have any spaces
	stackName := spartaCF.UserScopedStackName("SpartaXRay")
	var lambdaFunctions []*sparta.LambdaAWSInfo
	lambdaFunctions = append(lambdaFunctions, lambdaFn)
	err := spartaCGO.Main(stackName,
		stackName,
		lambdaFunctions,
		nil,
		nil)
	if err != nil {
		os.Exit(1)
	}
}

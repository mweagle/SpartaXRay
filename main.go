package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	gocf "github.com/mweagle/go-cloudformation"
)

// Standard AWS λ function
func helloXRay(event *json.RawMessage,
	context *sparta.LambdaContext,
	w http.ResponseWriter,
	logger *logrus.Logger) {
	fmt.Fprint(w, "Hello XRay World ☢️")
}

func xRaySampleFunctions() []*sparta.LambdaAWSInfo {
	// Default options for all lambda functions
	sampleFunctions := make([]*sparta.LambdaAWSInfo, 0)

	lambdaMaker := func(xRayIndex int) sparta.LambdaFunction {
		return func(event *json.RawMessage,
			context *sparta.LambdaContext,
			w http.ResponseWriter,
			logger *logrus.Logger) {
			fmt.Fprintf(w, "Hello World from XRay %d: ☢️", xRayIndex)
		}
	}

	// Create 6 functions
	for i := 0; i != 6; i++ {
		lambdaOptions := &sparta.LambdaFunctionOptions{
			MemorySize: 256,
			TracingConfig: &gocf.LambdaFunctionTracingConfig{
				Mode: gocf.String("Active"),
			},
			Tags: map[string]string{
				"myAccounting": "tag",
			},
			SpartaOptions: &sparta.SpartaOptions{
				Name: fmt.Sprintf("XRaySampleFunction%d", i),
			},
		}
		lambdaFn := sparta.NewLambda(sparta.IAMRoleDefinition{},
			lambdaMaker(i),
			lambdaOptions)
		sampleFunctions = append(sampleFunctions, lambdaFn)
	}
	return sampleFunctions
}

////////////////////////////////////////////////////////////////////////////////
// Main
func main() {
	// Setup some sample functions
	lambdaFunctions := xRaySampleFunctions()

	// Sanitize the name so that it doesn't have any spaces
	stackName := spartaCF.UserScopedStackName("SpartaXRay")

	// Setup the DashboardDecorator lambda hook
	workflowHooks := &sparta.WorkflowHooks{
		ServiceDecorator: sparta.DashboardDecorator(lambdaFunctions, 60),
	}

	err := sparta.MainEx(stackName,
		stackName,
		lambdaFunctions,
		nil,
		nil,
		workflowHooks,
		false)
	if err != nil {
		os.Exit(1)
	}
}

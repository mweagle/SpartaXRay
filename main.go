package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	sparta "github.com/mweagle/Sparta"
	spartaCF "github.com/mweagle/Sparta/aws/cloudformation"
	spartaDecorator "github.com/mweagle/Sparta/decorator"
	gocf "github.com/mweagle/go-cloudformation"
	"github.com/sirupsen/logrus"
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

	lambdaMaker := func(xRayIndex int) interface{} {
		return func(ctx context.Context) (string, error) {
			return fmt.Sprintf("Hello World from XRay %d: ☢️", xRayIndex), nil
		}
	}

	// Create 6 functions
	for i := 0; i != 6; i++ {
		lambdaOptions := &sparta.LambdaFunctionOptions{
			Timeout:    10,
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
		lambdaHandler := lambdaMaker(i)
		lambdaFn, lambdaFnErr := sparta.NewAWSLambda(fmt.Sprintf("XRayHandler%d", i),
			lambdaHandler,
			sparta.IAMRoleDefinition{})
		if lambdaFnErr != nil {
			panic(lambdaFnErr.Error())
		}
		lambdaFn.Options = lambdaOptions
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
		ServiceDecorators: []sparta.ServiceDecoratorHookHandler{
			spartaDecorator.DashboardDecorator(lambdaFunctions, 60),
		},
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

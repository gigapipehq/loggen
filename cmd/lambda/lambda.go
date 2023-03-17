package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	"github.com/gigapipehq/loggen/internal/cmd"
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/progress"
)

var lambdaCMD = &cobra.Command{
	Use:   "lambda",
	Short: "Run the generator in AWS Lambda mode",
	Run: func(_ *cobra.Command, _ []string) {
		lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			statusCode := fiber.StatusOK
			resp := "Ok"
			cfg := &config.Config{}
			if err := cfg.UnmarshalJSON([]byte(req.Body)); err != nil {
				statusCode = fiber.StatusBadRequest
				resp = err.Error()
			}

			if err := cmd.Do(ctx, cfg, progress.NewLambda()); err != nil {
				statusCode = fiber.StatusInternalServerError
				resp = err.Error()
			}
			return events.APIGatewayProxyResponse{
				StatusCode: statusCode,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: resp,
			}, nil
		})
	},
}

func init() {
	lambdaCMD.AddCommand()
}

func CMD() *cobra.Command {
	return lambdaCMD
}

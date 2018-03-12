package lambda

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	c.ListAliases(&lambda.ListAliasesInput{})

	// 		p.events.Emit("metrics.value", event.Fields{
	// 		"name":   s.Name,
	// 		"value":  s.Value(),
	// 		"memory": p.config.Lambda.Memory,
	// 	})

	return nil
}

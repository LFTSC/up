package lambda

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	var marker *string

	for {
		res, err := c.ListAliases(&lambda.ListAliasesInput{
			FunctionName: &p.config.Name,
			Marker:       marker,
			MaxItems:     aws.Int64(10000),
		})

		if err != nil {
			return errors.Wrap(err, "listing aliases")
		}

		{
			enc := json.NewEncoder(os.Stderr)
			enc.SetIndent("", "  ")
			enc.Encode(res.Aliases)
		}

		marker = res.NextMarker
		if marker == nil {
			break
		}
	}

	// list aliases
	// output git aliases
	// star them if matching prod etc

	// 		p.events.Emit("metrics.value", event.Fields{
	// 		"name":   s.Name,
	// 		"value":  s.Value(),
	// 		"memory": p.config.Lambda.Memory,
	// 	})

	return nil
}

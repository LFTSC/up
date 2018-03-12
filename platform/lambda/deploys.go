package lambda

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	aliases, err := getAliases(c, p.config.Name)
	if err != nil {
		return errors.Wrap(err, "fetching aliases")
	}

	{
		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "  ")
		enc.Encode(aliases)
	}

	// 		p.events.Emit("metrics.value", event.Fields{
	// 		"name":   s.Name,
	// 		"value":  s.Value(),
	// 		"memory": p.config.Lambda.Memory,
	// 	})

	return nil
}

// getAliases returns function aliases sorted by version.
func getAliases(c *lambda.Lambda, name string) (aliases []*lambda.AliasConfiguration, err error) {
	var marker *string

	for {
		res, err := c.ListAliases(&lambda.ListAliasesInput{
			FunctionName: &name,
			Marker:       marker,
			MaxItems:     aws.Int64(10000),
		})

		if err != nil {
			return nil, err
		}

		aliases = append(aliases, res.Aliases...)

		marker = res.NextMarker
		if marker == nil {
			break
		}
	}

	sort.Slice(aliases, func(i int, j int) bool {
		a := aliases[i]
		b := aliases[j]
		return *a.FunctionVersion > *b.FunctionVersion
	})

	return
}

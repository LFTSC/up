package lambda

import (
	"fmt"
	"sort"

	"github.com/apex/up/internal/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

// TODO: parallelize?

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	versions, err := getVersions(c, p.config.Name)
	if err != nil {
		return errors.Wrap(err, "fetching versions")
	}

	// aliases, err := getAliases(c, p.config.Name)
	// if err != nil {
	// 	return errors.Wrap(err, "fetching aliases")
	// }

	sortVersionsDesc(versions)
	defer util.Pad()()

	for _, v := range versions {
		version := *v.Version

		if version == "$LATEST" {
			continue
		}

		stage := *v.Environment.Variables["UP_STAGE"]
		commit := v.Environment.Variables["UP_COMMIT"]
		if commit != nil {
			fmt.Printf("  %s -> %s (%s)\n", stage, *commit, version)
		}
	}

	// p.events.Emit("platform.deploys", event.Fields{
	// 	"aliases": aliases,
	// })

	return nil
}

// getAliases returns all function aliases.
func getAliases(c *lambda.Lambda, name string) (aliases []*lambda.AliasConfiguration, err error) {
	var marker *string

	for {
		res, err := c.ListAliases(&lambda.ListAliasesInput{
			FunctionName: &name,
			Marker:       marker,
			MaxItems:     aws.Int64(5000),
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

	return
}

// getVersions returns all function versions.
func getVersions(c *lambda.Lambda, name string) (versions []*lambda.FunctionConfiguration, err error) {
	var marker *string

	for {
		res, err := c.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
			FunctionName: &name,
			MaxItems:     aws.Int64(5000),
		})

		if err != nil {
			return nil, err
		}

		versions = append(versions, res.Versions...)

		marker = res.NextMarker
		if marker == nil {
			break
		}
	}

	return
}

// sortVersionsDesc sorts versions descending.
func sortVersionsDesc(versions []*lambda.FunctionConfiguration) {
	sort.Slice(versions, func(i int, j int) bool {
		a := *versions[i].Version
		b := *versions[j].Version
		return a > b
	})
}

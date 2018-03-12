package lambda

import (
	"fmt"
	"sort"

	"github.com/apex/up/internal/util"
	"github.com/araddon/dateparse"
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

	for _, f := range versions {
		if *f.Version == "$LATEST" {
			continue
		}

		showFunction(f)
	}

	return nil
}

// showFunction outputs the change.
func showFunction(f *lambda.FunctionConfiguration) {
	commit := f.Environment.Variables["UP_COMMIT"]
	stage := *f.Environment.Variables["UP_STAGE"]
	created := dateparse.MustParse(*f.LastModified)
	version := *f.Version

	// no git commit
	if commit == nil || *commit == "" {
		fmt.Printf("  %15s -> %s %s\n", stage, version, created)
		return
	}

	// git commit
	fmt.Printf("  %15s -> %s (%s) %s\n", stage, *commit, version, created)
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

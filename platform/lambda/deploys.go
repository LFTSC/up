package lambda

import (
	"sort"
	"time"

	"github.com/apex/up/internal/table"
	"github.com/apex/up/internal/util"
	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

// ShowDeploys implementation.
func (p *Platform) ShowDeploys(region string) error {
	s := session.New(aws.NewConfig().WithRegion(region))
	c := lambda.New(s)

	versions, err := getVersions(c, p.config.Name)
	if err != nil {
		return errors.Wrap(err, "fetching versions")
	}

	sortVersionsDesc(versions)
	t := table.New()

	for _, f := range versions {
		if *f.Version == "$LATEST" {
			continue
		}

		addFunction(t, f)
	}

	defer util.Pad()()
	t.Println()

	return nil
}

// addFunction adds function to table.
func addFunction(t *table.Table, f *lambda.FunctionConfiguration) {
	commit := f.Environment.Variables["UP_COMMIT"]
	stage := *f.Environment.Variables["UP_STAGE"]
	created := dateparse.MustParse(*f.LastModified)
	date := formatDate(created)
	version := *f.Version

	t.AddRow(table.Row{
		{Text: stage},
		{Text: util.DefaultString(commit, "–")},
		{Text: version},
		{Text: date},
	})
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

// formatDate formats t relative to now.
func formatDate(t time.Time) string {
	switch d := time.Since(t); {
	case d <= 12*time.Hour:
		return humanize.RelTime(time.Now(), t, "from now", "ago")
	case d <= 24*time.Hour:
		return t.Format(`Today at 03:04:05pm`)
	default:
		return t.Format(`Jan 2` + util.DateSuffix(t) + ` 03:04:05pm`)
	}
}

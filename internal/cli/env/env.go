package env

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/tj/kingpin"

	"github.com/apex/up"
	"github.com/apex/up/internal/cli/root"
	"github.com/apex/up/internal/colors"
	"github.com/apex/up/internal/secret"
	"github.com/apex/up/internal/stats"
	"github.com/apex/up/internal/table"
	"github.com/apex/up/internal/util"
	"github.com/apex/up/internal/validate"
)

func init() {
	cmd := root.Command("env", "Manage encrypted env variables.")
	cmd.Example(`up env`, "List variables available to all stages.")
	cmd.Example(`up env -d`, "List decrypted values.")
	cmd.Example(`up env add MONGO_URL "mongodb://db1.example.net:2500/" -s production`, "Add a production env variable.")
	cmd.Example(`up env add MONGO_URL "mongodb://db2.example.net:2500/" -s staging`, "Add a staging env variable.")
	cmd.Example(`up env add S3_KEY xxxxxxx`, "Add add a global env variable for all stages.")
	cmd.Example(`up env add S3_KEY xxxxxxx -s production`, "Add a stage specific env var to override the previous.")
	cmd.Example(`up env add -c DB_USER tobi`, "Add a cleartext env var.")
	cmd.Example(`up env add -d 'Mongo password' DB_PASS xxxxxxx`, "Add a description.")
	cmd.Example(`up env get DB_URL`, "Get a variable value.")
	cmd.Example(`up env rm S3_KEY`, "Remove a variable.")
	cmd.Example(`up env rm S3_KEY -s production`, "Remove a production variable.")
	add(cmd)
	get(cmd)
	list(cmd)
	remove(cmd)
}

// get variables.
func get(cmd *kingpin.Cmd) {
	c := cmd.Command("get", "Get a variable value.")
	key := c.Arg("name", "Variable name.").Required().String()
	stage := c.Flag("stage", "Target stage name.").Short('s').Default("all").String()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stages := append(c.Stages.Names(), "all")
		if err := validate.List(*stage, stages); err != nil {
			return err
		}
		normalizeStage(stage)

		stats.Track("Get Secret", nil)

		v, err := p.Secrets(*stage).Get(*key)
		if err != nil {
			return errors.Wrap(err, "listing secrets")
		}

		fmt.Printf("%s", v)

		return nil
	})
}

// list variables.
func list(cmd *kingpin.Cmd) {
	c := cmd.Command("ls", "List variables.").Alias("list").Default()
	decrypt := c.Flag("decrypt", "Output decrypted values.").Short('d').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stats.Track("List Secrets", map[string]interface{}{
			"decrypt": *decrypt,
		})

		secrets, err := p.Secrets("").List(*decrypt)
		if err != nil {
			return errors.Wrap(err, "listing secrets")
		}

		if len(secrets) == 0 {
			util.LogPad("No environment variables defined. See `up help env` for examples.")
			return nil
		}

		grouped := secret.GroupByStage(secrets)
		t := table.New()

		stages := append([]string{"all"}, c.Stages.Names()...)
		for _, name := range stages {
			secrets, ok := grouped[name]
			if !ok {
				continue
			}

			t.AddRow(table.Row{
				{
					Text: colors.Bold(fmt.Sprintf("\n%s\n", strings.Title(name))),
					Span: 4,
				},
			})

			rows(t, secrets)
		}

		t.Println()
		println()

		return nil
	})
}

// add variables.
func add(cmd *kingpin.Cmd) {
	c := cmd.Command("add", "Add a variable.").Alias("set")
	key := c.Arg("name", "Variable name.").Required().String()
	val := c.Arg("value", "Variable value.").Required().String()
	stage := c.Flag("stage", "Target stage name.").Short('s').Default("all").String()
	desc := c.Flag("desc", "Variable description message.").Short('d').String()
	clear := c.Flag("clear", "Store as cleartext (unencrypted).").Short('c').Bool()

	c.Action(func(_ *kingpin.ParseContext) error {
		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stages := append(c.Stages.Names(), "all")
		if err := validate.List(*stage, stages); err != nil {
			return err
		}
		normalizeStage(stage)

		stats.Track("Add Secret", map[string]interface{}{
			"cleartext": *clear,
			"stage":     *stage,
			"has_desc":  *desc != "",
		})

		if err := p.Secrets(*stage).Add(*key, *val, *desc, *clear); err != nil {
			return errors.Wrap(err, "adding secret")
		}

		util.LogPad("Added " + *key)

		return nil
	})
}

// remove variables.
func remove(cmd *kingpin.Cmd) {
	c := cmd.Command("rm", "Remove a variable.").Alias("remove")
	stage := c.Flag("stage", "Target stage name.").Short('s').Default("all").String()
	key := c.Arg("name", "Variable name.").Required().String()

	c.Action(func(_ *kingpin.ParseContext) error {
		defer util.Pad()()

		c, p, err := root.Init()
		if err != nil {
			return errors.Wrap(err, "initializing")
		}

		stages := append(c.Stages.Names(), "all")
		if err := validate.List(*stage, stages); err != nil {
			return err
		}
		normalizeStage(stage)

		stats.Track("Remove Secret", nil)

		if err := p.Secrets(*stage).Remove(*key); err != nil {
			return errors.Wrap(err, "removing secret")
		}

		util.LogPad("Removed " + *key)

		return nil
	})
}

// rows helper.
func rows(t *table.Table, secrets []*up.Secret) {
	for _, s := range secrets {
		mod := fmt.Sprintf("Modified %s", humanize.Time(s.LastModified))
		if u := s.LastModifiedUser; u != "" {
			mod += fmt.Sprintf(" by %s", u)
		}
		desc := colors.Gray(util.DefaultString(&s.Description, "-"))
		val := colors.Gray(util.DefaultString(&s.Value, "-"))

		t.AddRow(table.Row{
			{
				Text: colors.Purple(s.Name),
			},
			{
				Text: val,
			},
			{
				Text: desc,
			},
			{
				Text: mod,
			},
		})
	}
}

// normalizeStage normalizes "all" which is internally represented as "".
func normalizeStage(s *string) {
	if *s == "all" {
		*s = ""
	}
}

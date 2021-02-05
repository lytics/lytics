package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

func init() {
	addCommand(cli.Command{
		Name:     "ml",
		Usage:    "Machine Learning Info",
		Category: "ML API",
		Subcommands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "Show details of requested segment",
				Action: mlGet,
			},
			{
				Name:   "list",
				Usage:  "List Machine Learning Models",
				Action: mlList,
			},
		},
	})
}
func mlGet(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetMLModel(id)
	exitIfErr(err, "could not get segment %q from API", id)
	resultWrite(c, &item, fmt.Sprintf("segment_%s", item.Name))
	return nil
}

func mlList(c *cli.Context) error {
	items, err := client.GetMLModels()
	exitIfErr(err, "could not get segment list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list, "segment_list")
	return nil
}

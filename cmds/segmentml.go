package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli"
)

var segMlOutput string

func init() {
	addCommand(cli.Command{
		Name:     "segmentml",
		Usage:    "SegmentML Info",
		Category: "Data API",
		Action:   run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Usage:       "Specify what segmentML CSV's or tables to output default: all; individual options: features, predictions, overview",
				Value:       "all",
				Destination: &segMlOutput,
			},
		},
	})
}

func run(c *cli.Context) error {
	var err error
	switch segMlOutput {
	case "features":
		err = segMlFeatures(c)
	case "predictions":
		err = segMlPredictions(c)
	case "overview":
		err = segMlOverview(c)
	case "all":
		err = segMlFeatures(c)
		err = segMlPredictions(c)
		err = segMlOverview(c)
	default:
		return fmt.Errorf("specify what segmentML table to output: all, features, predictions, or overview")
	}
	if err != nil {
		return fmt.Errorf("error creating table %v", err)
	}
	return nil
}

func segMlFeatures(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	SegML, err := client.GetSegmentMLModel(id)
	exitIfErr(err, "could not get segment list")
	list := make([]lytics.TableWriter, len(SegML.Features))
	for i, feat := range SegML.Features {
		list[i] = feat
	}
	if len(list) == 0 {
		return fmt.Errorf("no features")
	}
	name := fmt.Sprintf("Features-%s", SegML.Name)
	resultWrite(c, list, name)
	return nil
}

func segMlPredictions(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	SegML, err := client.GetSegmentMLModel(id)
	exitIfErr(err, "could not get segment list")
	predictions := SegML.GetPredictions()

	list := make([]lytics.TableWriter, len(predictions))
	for i, pred := range predictions {
		list[i] = pred
	}
	if len(list) == 0 {
		return fmt.Errorf("no predictions")
	}
	name := fmt.Sprintf("predictions-%s", SegML.Name)
	resultWrite(c, list, name)
	return nil
}

func segMlOverview(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	SegML, err := client.GetSegmentMLModel(id)
	exitIfErr(err, "could not get segment list")
	name := fmt.Sprintf("Overview-%s", SegML.Name)
	resultWrite(c, &SegML, name)
	return nil
}

package cmds

import (
	"fmt"

	lytics "github.com/lytics/go-lytics"
	"github.com/urfave/cli/v2"
)

func init() {
	addCommand(cli.Command{
		Name:     "segment",
		Usage:    "Segment Info",
		Category: "Management API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "table",
				Usage: "Table to limit list of segments",
				Value: "user",
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "Show details of requested segment",
				Action: segmentGet,
			},
			{
				Name:   "list",
				Usage:  "List Segments",
				Action: segmentList,
			},
			{
				Name:   "listql",
				Usage:  "List SegmentQL Queries",
				Action: segmentQlList,
			},
			{
				Name:      "scan",
				Usage:     "List Entities in a Segment.  NOTE, this is new-line delimitted json output.",
				ArgsUsage: "[id or slug of Segment]",
				Action:    segmentScan,
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "limit",
						Usage: "Limit to x entities in scan list",
						Value: 0,
					},
				},
			},
		},
	})
}
func segmentGet(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	item, err := client.GetSegment(id)
	exitIfErr(err, "could not get segment %q from API", id)
	resultWrite(c, &item, fmt.Sprintf("segment_%s", item.Name))
	return nil
}
func segmentList(c *cli.Context) error {
	items, err := client.GetSegments(c.String("table"))
	exitIfErr(err, "could not get segment list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = item
	}
	resultWrite(c, list, "segment_list")
	return nil
}
func segmentQlList(c *cli.Context) error {
	items, err := client.GetSegments(c.String("table"))
	exitIfErr(err, "could not get SegmentQL list")
	list := make([]lytics.TableWriter, len(items))
	for i, item := range items {
		list[i] = &segmentQl{item}
	}
	resultWrite(c, list, "SegmentQL list")
	return nil
}
func segmentScan(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("expected one arg (id)")
	}
	id := c.Args().First()
	limit := c.Int64("limit")
	getEntityScan(id, int(limit), func(e *lytics.Entity) {
		fmt.Println(e.PrettyJson())
	})
	return nil
}

type segmentQl struct {
	*lytics.Segment
}

func (m *segmentQl) Headers() []interface{} {
	return []interface{}{
		"ID", "alias", "ql",
	}
}
func (m *segmentQl) Row() []interface{} {
	return []interface{}{
		m.Id, m.SlugName, m.FilterQL,
	}
}

/*
func (c *Cli) getSegmentSizes(segments []string) (interface{}, error) {
	if len(segments) == 1 {
		segment, err := c.Client.GetSegmentSize(segments[0])
		if err != nil {
			return nil, err
		}

		return segment, nil
	} else {
		segments, err := c.Client.GetSegmentSizes(segments)
		if err != nil {
			return nil, err
		}

		return segments, nil
	}
}

func (c *Cli) getSegmentAttributions(segments []string, limit int) (interface{}, error) {
	attributions, err := c.Client.GetSegmentAttribution(segments)
	if err != nil {
		return nil, err
	}

	return attributions, nil
}
*/

func getEntityScan(segmentIDOrQl string, limit int, handler lytics.EntityHandler) {

	scan := client.PageSegment(segmentIDOrQl)
	ct := 0
	// handle processing the entities
	for {
		e := scan.Next()
		if e == nil {
			break
		}
		handler(e)
		ct++
		if limit > 0 && ct == limit {
			return
		}
	}
}

package command

import (
	"fmt"
	"strings"

	lytics "github.com/lytics/go-lytics"
)

func segmentCommands(api *apiCommand) map[string]*command {
	c := &segment{apiCommand: api}
	return map[string]*command{
		"list": &command{c.HelpList, c.List, "Segment List, segments for account."},
		"show": &command{c.HelpGet, c.Get, "Segment Show Summary."},
	}
}

type segment struct {
	*apiCommand
}

func (c *segment) HelpList() string {
	helpText := fmt.Sprintf(`
Usage: lytics segment list [options]

  List segments

%s
`, globalHelp)
	return strings.TrimSpace(helpText)
}
func (c *segment) HelpGet() string {
	helpText := fmt.Sprintf(`
Usage: lytics segment show [options] id

  Get Segment and show summary

%s

Options:
    --table="user"
`, globalHelp)
	return strings.TrimSpace(helpText)
}

func (c *segment) Get(args []string) int {

	c.init(args, c.HelpGet)
	id := c.f.Arg(0)
	if id == "" {
		c.ui.Error("Must provide segment ID/Slug")
		return 1
	}
	c.cols = []string{"name", "SegKind", "created", "filterql"}

	segment, err := c.l.GetSegment(id)
	c.exitIfErr(err, "Could not get segment")

	c.writeSingle(segment)
	return 0
}

func (c *segment) List(args []string) int {
	c.init(args, c.HelpList)
	table := c.f.Arg(0)
	if table == "" {
		table = "user"
	}

	c.cols = []string{"name", "SegKind", "created", "filterql"}

	segments, err := c.l.GetSegments(table)
	c.exitIfErr(err, "Could not get segment")
	items := make([]interface{}, len(segments))
	for i, u := range segments {
		items[i] = u
	}
	c.writeList(items)
	return 0
}

func (c *Cli) getSegments(table string, segments []string) (interface{}, error) {
	if len(segments) == 1 {
		segment, err := c.Client.GetSegment(segments[0])
		if err != nil {
			return nil, err
		}

		return segment, nil
	} else {
		segments, err := c.Client.GetSegments(table)
		if err != nil {
			return nil, err
		}

		return segments, nil
	}
}

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

func (c *Cli) getEntityScan(segmentIdOrQl string, limit int, handler lytics.EntityHandler) {

	scan := c.Client.PageSegment(segmentIdOrQl)

	ct := 0
	// handle processing the entities
	for {
		e := scan.Next()
		if e == nil {
			break
		}
		handler(&e)
		ct++
		if limit > 0 && ct == limit {
			return
		}
	}
}

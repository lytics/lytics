package main

import (
	lytics "github.com/lytics/go-lytics"
)

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

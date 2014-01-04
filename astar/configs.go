package astar

import (
    "fmt"
    "math"
)

var _ = fmt.Sprint()

//######################################################################
//######################################################################

type pointToPoint struct {
}

// Basic point to point routing, only a single source
// is supported and it will panic if given multiple sources
//
// Weights are calulated by summing the tiles fill_weight, the total distance traveled
// and the current distance from the target
func NewPointToPoint() AStarConfig {
    p2p := &pointToPoint{}

    return p2p
}

func (p2p *pointToPoint) SetWeight(p *PathPoint, fill_weight int, end []Point, end_map map[Point]bool) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }

    if fill_weight == -1 {
        return false
    }

    p.Weight = p.FillWeight + p.DistTraveled + p.Point.Dist(end[0])

    return true
}

func (p2p *pointToPoint) IsEnd(p Point, end []Point, end_map map[Point]bool) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }
    return p == end[0]
}

//######################################################################
//######################################################################

type rowToRow struct {
}

// Based off the PointToPoint config except that it uses row based targeting.
// The column value is ignored when calculating the weight and when determining
// if we have reached the end.

// A single point should be given for the source which will determine the starting row.
// for the target you should provide every valid entry on the target row for the best results.
// you do not have to but the path may look a little strange sometimes.
func NewRowToRow() AStarConfig {
    r2r := &rowToRow{}
    return r2r
}

func (r2r *rowToRow) SetWeight(p *PathPoint, fill_weight int, end []Point, end_map map[Point]bool) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }

    if fill_weight == -1 {
        return false
    }

    p.Weight = p.FillWeight + p.DistTraveled + int(math.Abs(float64(p.Row-end[0].Row)))

    return true
}

func (r2r *rowToRow) IsEnd(p Point, end []Point, end_map map[Point]bool) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }
    return p.Row == end[0].Row
}

//######################################################################
//######################################################################

type pointToList struct {
}

// point to list routing, from any single point to a list of points. for example
// to find your way back to the main path.
//
// Both multiple targets and sources are supported but is slower than the others
//
// Weights are calulated by summing the tiles fill_weight, the total distance traveled
// and the current distance from the closeset target
func NewPointToList() AStarConfig {
    p2l := &pointToList{}

    return p2l
}

func (p2l *pointToList) SetWeight(p *PathPoint, fill_weight int, end []Point, end_map map[Point]bool) bool {
    if fill_weight == -1 {
        return false
    }

    min_dist := -1
    for _, end_p := range end {
        dist := p.Point.Dist(end_p)
        if min_dist == -1 || dist < min_dist {
            min_dist = dist
        }
    }

    p.Weight = p.FillWeight + p.DistTraveled + min_dist

    return true
}

func (p2l *pointToList) IsEnd(p Point, end []Point, end_map map[Point]bool) bool {
    return end_map[p]
}

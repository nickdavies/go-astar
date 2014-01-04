package astar

import (
    "fmt"
    "math"
)

var _ = fmt.Sprint()

//######################################################################
//######################################################################

type pointToPoint struct {
    *AStarBaseStruct
}

// Basic point to point routing, only a single source
// is supported and it will panic if given multiple sources
//
// Weights are calulated by summing the tiles fill_weight, the total distance traveled
// and the current distance from the target
func NewPointToPoint(rows, cols int) AStar {
    p2p := &pointToPoint{
        AStarBaseStruct: NewAStarBaseStruct(rows, cols),
    }

    p2p.AStarBaseStruct.Config = p2p

    return p2p
}

func (p2p *pointToPoint) SetWeight(p *PathPoint, fill_weight int, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }

    if fill_weight == -1 {
        return false
    }

    p.Weight = p.FillWeight + p.DistTraveled + p.Point.Dist(end[0])

    return true
}

func (p2p *pointToPoint) IsEnd(p Point, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }
    return p == end[0]
}

//######################################################################
//######################################################################

type rowToRow struct {
    *AStarBaseStruct
}

// Based off the PointToPoint config except that it uses row based targeting.
// The column value is ignored when calculating the weight and when determining
// if we have reached the end.

// A single point should be given for the source which will determine the starting row.
// for the target you should provide every valid entry on the target row for the best results.
// you do not have to but the path may look a little strange sometimes.
func NewRowToRow(rows, cols int) AStar {
    r2r := &rowToRow{
        AStarBaseStruct: NewAStarBaseStruct(rows, cols),
    }

    r2r.AStarBaseStruct.Config = r2r

    return r2r
}

func (r2r *rowToRow) SetWeight(p *PathPoint, fill_weight int, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }

    if fill_weight == -1 {
        return false
    }

    p.Weight = p.FillWeight + p.DistTraveled + int(math.Abs(float64(p.Row-end[0].Row)))

    return true
}

func (r2r *rowToRow) IsEnd(p Point, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }
    return p.Row == end[0].Row
}

package astar

import "fmt"

var _ = fmt.Sprint()

type PointToPoint struct {
    *AStarBaseStruct
}

func NewPointToPoint(rows, cols int) AStar {
    p2p := &PointToPoint{
        AStarBaseStruct: NewAStarBaseStruct(rows, cols),
    }

    p2p.AStarBaseStruct.Config = p2p

    return p2p
}

func (p2p *PointToPoint) SetWeight(p *PathPoint, fill_weight int, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }

    if fill_weight == -1 {
        return false
    }

    p.Weight = p.FillWeight + p.DistTraveled + p.Point.Dist(end[0])

    return true
}

func (p2p *PointToPoint) IsEnd(p Point, end []Point) bool {
    if len(end) != 1 {
        panic("Invalid end specified")
    }
    return p == end[0]
}

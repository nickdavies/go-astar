package astar

import "fmt"

var _ = fmt.Sprint()

type PointToPoint struct {
    *AStarBase

    target Point
}

func NewPointToPoint (rows, cols int) AStar {
    p2p := &PointToPoint{
        AStarBase: NewAStarBase(rows, cols),
    }

    p2p.AStarBase.config = p2p

    return p2p
}

func (p2p *PointToPoint) SetWeight(p *PathPoint, fill_weight int, target []Point) bool {
    if len(target) != 1 {
        panic("Invalid Target Specified")
    }

    p.Weight = p.FillWeight + p.DistTraveled + p.Point.Dist(target[0])

    return true
}

func (p2p *PointToPoint) IsTarget(p Point, target []Point) bool {
    if len(target) != 1 {
        panic("Invalid Target Specified")
    }
    return p == target[0]
}

package astar

import (
    "math"
)

type AStar interface {
    aStarBase
    AStarConfig
}

type aStarBase interface {
    FillTile(p Point, weight int)
    ClearTile(p Point)

    FindPath(source, target []Point) (*PathPoint, map[Point]*PathPoint)
}

type AStarConfig interface {
    IsTarget(p Point, target []Point) bool

    SetWeight(p *PathPoint, fill_weight int, target []Point) (allowed bool)
}

type AStarBase struct {
    config AStarConfig
    // A list of filled tiles and their weight
    filledTiles map[Point]int

    rows int
    cols int
}

func NewAStarBase(rows, cols int) *AStarBase {
    b := &AStarBase {
        rows: rows,
        cols: cols,

        filledTiles: make(map[Point]int),
    }
    var _ aStarBase = b

    return b
}

func (a *AStarBase) FillTile(p Point, weight int) {
    a.filledTiles[p] = weight
}

func (a *AStarBase) ClearTile(p Point) {
    delete(a.filledTiles, p)
}

func (a *AStarBase) FindPath(source, target []Point) (*PathPoint, map[Point]*PathPoint) {
    var openList = make(map[Point]*PathPoint)
    var closeList = make(map[Point]*PathPoint)

    for _, p := range source {
        openList[p] = &PathPoint{
            Point: p,
            Parent: nil,
        }
    }

    var current *PathPoint
    for {
        current = a.getMinWeight(openList)
        if current == nil || a.config.IsTarget(current.Point, target) {
            break
        }

        delete(openList, current.Point)
        closeList[current.Point] = current

        surrounding := a.getSurrounding(current.Point)

        for _, p := range surrounding {
            _, ok := closeList[p]
            if ok {
                continue
            }

            fill_weight := a.filledTiles[p]
            if fill_weight == -1 {
                continue
            }

            path_point := &PathPoint{
                Point: p,
                Parent: current,
                FillWeight: current.FillWeight + fill_weight,
                DistTraveled: current.DistTraveled + 1,
            }
            allowed := a.config.SetWeight(path_point, fill_weight, target)
            if !allowed {
                continue
            }

            existing_point, ok := openList[p]
            if !ok {
                openList[p] = path_point
            } else {
                if path_point.Weight < existing_point.Weight {
                    existing_point.Parent = path_point.Parent
                }
            }
        }
    }

    return current, closeList
}

func (a *AStarBase) getMinWeight(openList map[Point]*PathPoint) *PathPoint {
    var min *PathPoint = nil
    var minWeight int = 0

    for _, p := range openList {
        if min == nil || p.Weight < minWeight {
            min = p
            minWeight = p.Weight
        }
    }
    return min
}

func (a *AStarBase) getSurrounding(p Point) []Point {
    var surrounding []Point

    row, col := p.Row, p.Col

    // Can only look up if not at the top
    if row > 0 {
        surrounding = append(surrounding, Point{row - 1, col})
    }
    if row < a.rows {
        surrounding = append(surrounding, Point{row + 1, col})
    }

    if col > 0 {
        surrounding = append(surrounding, Point{row, col - 1})
    }
    if col < a.cols {
        surrounding = append(surrounding, Point{row, col + 1})
    }

    return surrounding
}

type Point struct {
    Row int
    Col int
}

type PathPoint struct {
    Point
    Parent *PathPoint

    Weight int
    FillWeight int
    DistTraveled int

    WeightData interface{}
}

func (p Point) Dist(other Point) int {
    return int(math.Abs(float64(p.Row - other.Row)) + math.Abs(float64(p.Col - other.Col)))
}


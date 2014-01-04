package astar

import (
    "math"
    "sync"
)

type AStar interface {
    AStarBase
    AStarConfig
}

type AStarBase interface {
    FillTile(p Point, weight int)
    ClearTile(p Point)

    FindPath(source, target []Point) *PathPoint
}

type AStarConfig interface {
    IsEnd(p Point, end []Point) bool

    SetWeight(p *PathPoint, fill_weight int, end []Point) (allowed bool)
}

type AStarBaseStruct struct {
    Config AStarConfig
    // A list of filled tiles and their weight
    tileLock    sync.Mutex
    filledTiles map[Point]int

    rows int
    cols int
}

func NewAStarBaseStruct(rows, cols int) *AStarBaseStruct {
    b := &AStarBaseStruct{
        rows: rows,
        cols: cols,

        filledTiles: make(map[Point]int),
    }
    var _ AStarBase = b

    return b
}

func (a *AStarBaseStruct) FillTile(p Point, weight int) {
    a.tileLock.Lock()
    defer a.tileLock.Unlock()

    a.filledTiles[p] = weight
}

func (a *AStarBaseStruct) ClearTile(p Point) {
    a.tileLock.Lock()
    defer a.tileLock.Unlock()

    delete(a.filledTiles, p)
}

func (a *AStarBaseStruct) FindPath(source, target []Point) *PathPoint {
    var openList = make(map[Point]*PathPoint)
    var closeList = make(map[Point]*PathPoint)

    for _, p := range target {
        openList[p] = &PathPoint{
            Point:  p,
            Parent: nil,
        }
    }

    var current *PathPoint
    for {
        current = a.getMinWeight(openList)

        a.tileLock.Lock()
        if current == nil || a.Config.IsEnd(current.Point, source) {
            a.tileLock.Unlock()
            break
        }
        a.tileLock.Unlock()

        delete(openList, current.Point)
        closeList[current.Point] = current

        surrounding := a.getSurrounding(current.Point)

        for _, p := range surrounding {
            _, ok := closeList[p]
            if ok {
                continue
            }

            a.tileLock.Lock()
            fill_weight := a.filledTiles[p]
            a.tileLock.Unlock()

            path_point := &PathPoint{
                Point:        p,
                Parent:       current,
                FillWeight:   current.FillWeight + fill_weight,
                DistTraveled: current.DistTraveled + 1,
            }

            a.tileLock.Lock()
            allowed := a.Config.SetWeight(path_point, fill_weight, source)
            a.tileLock.Unlock()

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

    return current
}

func (a *AStarBaseStruct) getMinWeight(openList map[Point]*PathPoint) *PathPoint {
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

func (a *AStarBaseStruct) getSurrounding(p Point) []Point {
    var surrounding []Point

    row, col := p.Row, p.Col

    // Can only look up if not at the top
    if row > 0 {
        surrounding = append(surrounding, Point{row - 1, col})
    }
    if row < a.rows-1 {
        surrounding = append(surrounding, Point{row + 1, col})
    }

    if col > 0 {
        surrounding = append(surrounding, Point{row, col - 1})
    }
    if col < a.cols-1 {
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

    Weight       int
    FillWeight   int
    DistTraveled int

    WeightData interface{}
}

func (p Point) Dist(other Point) int {
    return int(math.Abs(float64(p.Row-other.Row)) + math.Abs(float64(p.Col-other.Col)))
}

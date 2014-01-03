package astar

import (
    "math"
)

type AStar interface {
    FillTile(p Point, weight int)
    ClearTile(p Point)

    FindPath(source []Point, target Point, weightCalc WeightCalculation)
}

type Point struct {
    Row int
    Col int
}

type PathPoint struct {
    Point
    Parent *PathPoint

    Weight int
    DistTraveled int

    WeightData interface{}
}

func (p Point) Dist(other Point) int {
    return int(math.Abs(float64(p.Row - other.Row)) + math.Abs(float64(p.Col - other.Col)))
}

func NewAStar(Rows, Cols int) AStar {
    return &aStarStruct {
        rows: Rows,
        cols: Cols,

        filledTiles: make(map[Point]int),
    }
}

type aStarStruct struct {
    // A list of filled tiles and their weight
    filledTiles map[Point]int

    rows int
    cols int
}

func (a *aStarStruct) FillTile(p Point, weight int) {
    a.filledTiles[p] = weight
}

func (a *aStarStruct) ClearTile(p Point) {
    delete(a.filledTiles, p)
}

func (a *aStarStruct) FindPath(source []Point, target Point, weightCalc WeightCalculation) {
    var openList map[Point]*PathPoint
    var closeList map[Point]*PathPoint

    for _, p := range source {
        openList[p] = &PathPoint{
            Point: p,
            Parent: nil,
        }
    }

    var current *PathPoint
    for {
        current = a.getMinWeight(openList)
        if current == nil || current.Point == target {
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
            }
            weightCalc(path_point, fill_weight, target)

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

}

func (a *aStarStruct) getMinWeight(openList map[Point]*PathPoint) *PathPoint {
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

func (a *aStarStruct) getSurrounding(p Point) []Point {
    var surrounding []Point

    row, col := p.Row, p.Col

    // Can only look up if not at the top
    if row > 0 {
        surrounding = append(surrounding, Point{row - 1, col})
    }
    if row < a.rows {
        surrounding = append(surrounding, Point{row - 1, col})
    }

    if col > 0 {
        surrounding = append(surrounding, Point{row, col - 1})
    }
    if col < a.cols {
        surrounding = append(surrounding, Point{row, col - 1})
    }

    return surrounding
}



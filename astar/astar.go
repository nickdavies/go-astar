package astar

import (
    "math"
    "sync"
)

type AStar interface {
    // Fill a given tile with a given weight this is used for making certain areas more complicated
    // to cross than others. For example you may have a higher weight for a wall or mountain.
    // This weight will be given back to you in the SetWeight function
    // Inbuilt A*'s use -1 to determine that it can not be passed at all.
    FillTile(p Point, weight int)

    // Resets the weight back to 0 for a given tile
    ClearTile(p Point)

    // Calculate the easiest path from ANY element in source to ANY element in target.
    // There is no hard rules about which element will become the start and end (unless your config
    // enforces it).
    // The start of the path is returned to you. If no path exists then the function will
    // return nil as the path.
    FindPath(config AStarConfig, source, target []Point) *PathPoint
}

// The user built configuration that determines how weights are calculated and
// also determines the stopping condition
type AStarConfig interface {
    // Determine if a valid end point has been reached. The end parameter
    // is the value passed in as source because the algorithm works backwards.
    IsEnd(p Point, end []Point, end_map map[Point]bool) bool

    // Calculate and set the weight for p.
    // fill_weight is the weight assigned to the tile when FillTile was called
    // or 0 if it was never called for that tile.
    // end is also provided so you can perform calculations such as distance remaining.
    SetWeight(p *PathPoint, fill_weight int, end []Point, end_map map[Point]bool) (allowed bool)

    // PostProcess the path after it has been calculated this might be useful if you want do do things
    // like reverse it or add constant moves at the start or end etc.
    PostProcess(p *PathPoint, rows, cols int, filledTiles map[Point]int) (*PathPoint)
}

type gridStruct struct {
    // A list of filled tiles and their weight
    tileLock    sync.Mutex
    filledTiles map[Point]int

    rows int
    cols int
}

func NewAStar(rows, cols int) AStar {
    return &gridStruct{
        rows: rows,
        cols: cols,

        filledTiles: make(map[Point]int),
    }
}

func (a *gridStruct) FillTile(p Point, weight int) {
    a.tileLock.Lock()
    defer a.tileLock.Unlock()

    a.filledTiles[p] = weight
}

func (a *gridStruct) ClearTile(p Point) {
    a.tileLock.Lock()
    defer a.tileLock.Unlock()

    delete(a.filledTiles, p)
}

func (a *gridStruct) FindPath(config AStarConfig, source, target []Point) *PathPoint {
    var openList = make(map[Point]*PathPoint)
    var closeList = make(map[Point]*PathPoint)

    source_map := make(map[Point]bool)
    for _, p := range source {
        source_map[p] = true
    }

    a.tileLock.Lock()
    for _, p := range target {
        fill_weight := a.filledTiles[p]
        path_point := &PathPoint{
            Point:        p,
            Parent:       nil,
            DistTraveled: 0,
            FillWeight:   fill_weight,
        }

        allowed := config.SetWeight(path_point, fill_weight, source, source_map)
        if allowed {
            openList[p] = path_point
        }
    }

    a.tileLock.Unlock()

    var current *PathPoint
    for {
        current = a.getMinWeight(openList)

        a.tileLock.Lock()
        if current == nil || config.IsEnd(current.Point, source, source_map) {
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
            allowed := config.SetWeight(path_point, fill_weight, source, source_map)
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

    a.tileLock.Lock()
    current = config.PostProcess(current, a.rows, a.cols, a.filledTiles)
    a.tileLock.Unlock()

    return current
}

func (a *gridStruct) getMinWeight(openList map[Point]*PathPoint) *PathPoint {
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

func (a *gridStruct) getSurrounding(p Point) []Point {
    var surrounding []Point

    row, col := p.Row, p.Col

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

// A point along a path.
// FillWeight is the sum of all the fill weights so far and
// DistTraveled is the total distance traveled so far
//
// WeightData is an interface that can be set to anything that Config wants
// it will never be touched by the rest of the code but if you wish to
// have any custom data held per node you can use WeightData
type PathPoint struct {
    Point
    Parent *PathPoint

    Weight       int
    FillWeight   int
    DistTraveled int

    WeightData interface{}
}

// Manhattan distance NOT euclidean distance because in our routing we cant go diagonally between the points.
func (p Point) Dist(other Point) int {
    return int(math.Abs(float64(p.Row-other.Row)) + math.Abs(float64(p.Col-other.Col)))
}

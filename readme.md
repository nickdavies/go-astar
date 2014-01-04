# go-astar - Glang A* #

golang A-star library with plug-able weight and targeting.

go get: `go get github.com/nickdavies/go-astar/astar`

Godoc: http://godoc.org/github.com/nickdavies/go-astar/astar

### Simple Routing ###

Using simple point A to point B routing

```go
package main

import "fmt"

import "github.com/nickdavies/go-astar/astar"

func main () {
    rows := 3
    cols := 3

    // Build AStar object from existing
    // PointToPoint configuration
    a := astar.NewPointToPoint(rows, cols)

    // Make an invincible obsticle at (1,1)
    a.FillTile(astar.Point{1, 1}, -1) 

    // Path from one corner to the other
    source := []astar.Point{astar.Point{0,0}}
    target := []astar.Point{astar.Point{2,2}}

    path, _ := a.FindPath(source, target)

    for path != nil {
        fmt.Printf("At (%d, %d)\n", path.Col, path.Row)
        path = path.Parent
    }
}
```

### Custom Routing Logic ###

To use your own routing logic make a struct that has a `AStarBaseStruct` struct within it and
initialise the `AStarBaseStruct` with the `NewAStarBase` function and set `AStarBaseStruct.config = x` where x is
your struct. Then make sure your struct implements `AStarConfig` and you can use it to route.

For example if you wanted routing that would ignore walls something that ignores walls:

```go
type MyAStar struct {
    *astar.AStarBaseStruct
}

func NewMyAStar (rows, cols int) astar.AStar {
    my := &MyAStar{
        AStarBaseStruct: astar.NewAStarBase(rows, cols),
    }   
    my.AStarBaseStruct.Config = my
    return my
}

func (my *MyAStar) SetWeight(p *astar.PathPoint, fill_weight int, end []astar.Point) bool {
    p.Weight = p.DistTraveled + p.Point.Dist(end[0])
    return true
}

func (my *MyAStar) IsEnd(p astar.Point, end []astar.Point) bool {
    return p == end[0]
}
```

## Full Example ###

A full example with a proper grid and walls etc can be found in `example.go` in the root of the repo

## Thread Safe ###

__All operations are thread-safe__ so you can perform multiple path searches on the same grid at the same time.
A lock on the tiles is acquired before calling any of the Config functions.

The only exception to this is if you change the `Config` property from outside the struct which you should never do anyway.

## Contributing ##

MIT licenced.

Feel free to send me pull requests or report issues.

For any pull requests please go fmt your code with:

```
gofmt -l -w -tabs=false -tabwidth=4
```


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
    a := astar.NewAStar(rows, cols)
    p2p := astar.NewPointToPoint()

    // Make an invincible obsticle at (1,1)
    a.FillTile(astar.Point{1, 1}, -1) 

    // Path from one corner to the other
    source := []astar.Point{astar.Point{0,0}}
    target := []astar.Point{astar.Point{2,2}}

    path := a.FindPath(p2p, source, target)

    for path != nil {
        fmt.Printf("At (%d, %d)\n", path.Col, path.Row)
        path = path.Parent
    }
}
```

### Custom Routing Logic ###

To use your own routing logic make a struct that implements `AStarConfig` and then pass that
struct into your call to FindPath

You must also decide if you want any post processing to be applied to your path before it is
returned, if you want none then you should embed the `VoidPostProcess` struct into your struct or
if you want to reverse the path use the `ReversePostProcess` struct or define your own.

For example if you wanted routing that would ignore walls something that ignores walls:

```go
type MyAStar struct {
    VoidPostProcess
}

func NewMyAStar () astar.AStar {
    my := &MyAStar{
    }
    return my
}

func (my *MyAStar) SetWeight(p *astar.PathPoint, fill_weight int, end []astar.Point, end_map map[astar.Point]bool) bool {
    p.Weight = p.DistTraveled + p.Point.Dist(end[0])
    return true
}

func (my *MyAStar) IsEnd(p astar.Point, end []astar.Point, end_map map[astar.Point]bool) bool {
    return end_map[p]
}
```

## Full Example ###

A full example with a proper grid and walls etc can be found in `example.go` in the root of the repo

An interesting seed is 1388819613980807431 (you can set this at the top of the file or leave it for random)

## Thread Safe ###

__All operations are thread-safe__ so you can perform multiple path searches on the same grid at the same time.
A lock on the tiles is acquired before calling any of the Config functions.

## Contributing ##

MIT licenced.

Feel free to send me pull requests or report issues.

For any pull requests please go fmt your code with:

```
gofmt -l -w -tabs=false -tabwidth=4
```


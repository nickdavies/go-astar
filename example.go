package main

import (
    "fmt"
    "math/rand"
    "time"
)

import (
    "github.com/nickdavies/go-astar/astar"
)

func main() {

    var start_t int64
    var end_t int64

    var seed int64 = 0

    // Setup the aStar structs
    ast := astar.NewAStar(50, 50)

    p2p := astar.NewPointToPoint()
    p2l := astar.NewListToPoint(true)

    // Generate a random map
    grid, source, target, me := GenerateRandomMap(ast, seed, 50, 600, 24, 100000)
    PrintGrid(grid)

    // Route from source to target (point to point)
    start_t = time.Now().UnixNano()
    end := ast.FindPath(p2p, source, target)
    end_t = time.Now().UnixNano()

    first_path_t := float64(end_t-start_t) / float64(time.Millisecond)

    DrawPath(grid, end, "*")
    PrintGrid(grid)

    // record path as array so it can be used in the next search
    p := end
    path := make([]astar.Point, 0)
    for p != nil {
        path = append(path, p.Point)
        p = p.Parent
    }

    start_t = time.Now().UnixNano()
    end = ast.FindPath(p2l, path, me)
    end_t = time.Now().UnixNano()
    second_path_t := float64(end_t-start_t) / float64(time.Millisecond)

    DrawPath(grid, end, ".")
    PrintGrid(grid)

    fmt.Println("me", me)
    fmt.Println("end", end)
    fmt.Println("end_grid", grid[end.Row][end.Col])
    fmt.Println(first_path_t)
    fmt.Println(second_path_t)
}

func GenerateRandomMap(ast astar.AStar, map_seed int64, grid_size, wall_count, wall_size, wall_weight int) ([][]string, []astar.Point, []astar.Point, []astar.Point) {

    if map_seed == 0 {
        map_seed = time.Now().UnixNano()
    }

    fmt.Println("Map Seed", map_seed)
    rand.Seed(map_seed)

    grid := make([][]string, grid_size)
    for i := 0; i < len(grid); i++ {
        grid[i] = make([]string, grid_size)
    }

    for walls := 0; walls < wall_count; {
        size := GetRandInt(wall_size)
        direction := GetRandInt(2)

        if direction == 0 {
            c := GetRandInt(grid_size)
            r := GetRandInt(grid_size - size)

            for i := 0; i < size; i++ {
                grid[r+i][c] = "#"
                ast.FillTile(astar.Point{r + i, c}, wall_weight)
            }
        } else {
            c := GetRandInt(grid_size - size)
            r := GetRandInt(grid_size)

            for i := 0; i < size; i++ {
                grid[r][c+i] = "#"
                ast.FillTile(astar.Point{r, c + i}, wall_weight)
            }
        }
        walls += size
    }

    for i := 0; i < grid_size; i++ {
        grid[0][i] = "#"
        ast.FillTile(astar.Point{0, i}, -1)

        grid[i][0] = "#"
        ast.FillTile(astar.Point{i, 0}, -1)

        grid[grid_size-1][i] = "#"
        ast.FillTile(astar.Point{grid_size - 1, i}, -1)

        grid[i][grid_size-1] = "#"
        ast.FillTile(astar.Point{i, grid_size - 1}, -1)
    }

    source := make([]astar.Point, 1)
    for {
        r := GetRandInt(grid_size)
        c := GetRandInt(grid_size)

        if grid[r][c] != "#" {
            grid[r][c] = "a"

            source[0].Row = r
            source[0].Col = c
            break
        }
    }

    target := make([]astar.Point, 1)
    for {
        r := GetRandInt(grid_size)
        c := GetRandInt(grid_size)

        if grid[r][c] != "#" && grid[r][c] != "a" {
            grid[r][c] = "b"

            target[0].Row = r
            target[0].Col = c
            break
        }
    }

    me := make([]astar.Point, 1)
    for {
        r := GetRandInt(grid_size)
        c := GetRandInt(grid_size)

        if grid[r][c] != "#" && grid[r][c] != "a" && grid[r][c] != "b" {
            grid[r][c] = "c"

            me[0].Row = r
            me[0].Col = c
            break
        }
    }

    return grid, source, target, me
}

func DrawPath(grid [][]string, path *astar.PathPoint, path_char string) {
    for {
        if grid[path.Row][path.Col] == "#" {
            grid[path.Row][path.Col] = "X"
        } else if grid[path.Row][path.Col] == "" {
            grid[path.Row][path.Col] = path_char
        }

        path = path.Parent
        if path == nil {
            break
        }
    }
}

func PrintGrid(grid [][]string) {
    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[0]); j++ {
            if grid[i][j] == "" {
                fmt.Printf(" ")
            } else {
                fmt.Print(grid[i][j])
            }
        }
        fmt.Print("\n")
    }
    fmt.Print("\n")
}

func GetRandInt(limit int) int {
    return rand.Intn(limit)
}

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

    var seed int64 = 0

    grid, ast, source, target := GenerateRandomMap(seed, 50, 600, 24, 100000)

    PrintGrid(grid)
    end := ast.FindPath(source, target)

    DrawPath(grid, end, source, target)
    PrintGrid(grid)
    fmt.Println(end)
}

func GenerateRandomMap(map_seed int64, grid_size, wall_count, wall_size, wall_weight int) ([][]string, astar.AStar, []astar.Point, []astar.Point) {

    if map_seed == 0 {
        map_seed = time.Now().UnixNano()
    }

    fmt.Println("Map Seed", map_seed)
    rand.Seed(map_seed)

    //ast := astar.NewPointToPoint(grid_size, grid_size)
    ast := astar.NewRowToRow(grid_size, grid_size)

    grid := make([][]string, grid_size)
    for i := 0; i < len(grid); i++ {
        grid[i] = make([]string, grid_size)
    }

    for walls := 0; walls < wall_count; {
        size := GetRandInt(wall_size)
        direction := GetRandInt(1) + 1

        if direction == 0 {
            c := GetRandInt(grid_size)
            r := GetRandInt(grid_size-size-2) + 1

            for i := 0; i < size; i++ {
                grid[r+i][c] = "#"
                ast.FillTile(astar.Point{r + i, c}, wall_weight)
            }
        } else {
            c := GetRandInt(grid_size - size)
            r := GetRandInt(grid_size-2) + 1

            for i := 0; i < size; i++ {
                grid[r][c+i] = "#"
                ast.FillTile(astar.Point{r, c + i}, wall_weight)
            }
        }
        walls += size
    }

    for i := 0; i < grid_size; i++ {
        grid[i][0] = "#"
        ast.FillTile(astar.Point{i, 0}, -1)

        grid[i][grid_size-1] = "#"
        ast.FillTile(astar.Point{i, grid_size - 1}, -1)
    }

    source := make([]astar.Point, 1)

    r := 0
    for i := 0; i < grid_size; i++ {
        grid[r][i] = "A"
    }
    source[0].Row = r

    target := make([]astar.Point, grid_size)
    r = grid_size - 1
    for i := 0; i < grid_size; i++ {
        grid[r][i] = "B"
        target[i].Row = r
        target[i].Col = i
    }

    return grid, ast, source, target
}

func DrawPath(grid [][]string, path *astar.PathPoint, source, target []astar.Point) {
    for {
        if path.Row == source[0].Row && path.Col == source[0].Col {
            grid[path.Row][path.Col] = "A"
        } else if path.Row == target[0].Row && path.Col == target[0].Col {
            grid[path.Row][path.Col] = "B"
        } else {
            if grid[path.Row][path.Col] == "#" {
                grid[path.Row][path.Col] = "X"
            } else {
                grid[path.Row][path.Col] = "*"
            }
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
}

func GetRandInt(limit int) int {
    return rand.Intn(limit)
}

package main

import (
    "time"
    "fmt"
    "math/rand"
)

import (
    "github.com/nickdavies/go-astar/astar"
)

func main() {

    grid, ast, source, target := GenerateRandomMap(0, 50, 600, 24, 100000)

    PrintGrid(grid)
    end, _ := ast.FindPath(source, target)

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

    ast := astar.NewPointToPoint(grid_size, grid_size)

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
                grid[r + i][c] = "#"
                ast.FillTile(astar.Point{r + i, c}, wall_weight)
            }
        } else {
            c := GetRandInt(grid_size - size)
            r := GetRandInt(grid_size)

            for i := 0; i < size; i++ {
                grid[r][c + i] = "#"
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

        grid[grid_size - 1][i] = "#"
        ast.FillTile(astar.Point{grid_size - 1, i}, -1)

        grid[i][grid_size - 1] = "#"
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

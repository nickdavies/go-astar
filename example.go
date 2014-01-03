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

    tower_weight := 100000
    seed := time.Now().UnixNano()

    fmt.Println(seed)
    rand.Seed(seed)

    ast := astar.NewAStar(50, 50)

    grid := make([][]string, 50)
    for i := 0; i < len(grid); i++ {
        grid[i] = make([]string, 50)
    }

    for wall_count := 0; wall_count < 600; {
        size := GetRandInt(24)
        direction := GetRandInt(2)

        if direction == 0 {
            c := GetRandInt(50)
            r := GetRandInt(50 - size)

            for i := 0; i < size; i++ {
                grid[r + i][c] = "#"
                ast.FillTile(astar.Point{r + i, c}, tower_weight)
            }
        } else {
            c := GetRandInt(50 - size)
            r := GetRandInt(50)

            for i := 0; i < size; i++ {
                grid[r][c + i] = "#"
                ast.FillTile(astar.Point{r, c + i}, tower_weight)
            }
        }
        wall_count += size
    }

    for i := 0; i < 50; i++ {
        grid[0][i] = "#"
        ast.FillTile(astar.Point{0, i}, -1)

        grid[i][0] = "#"
        ast.FillTile(astar.Point{i, 0}, -1)

        grid[49][i] = "#"
        ast.FillTile(astar.Point{49, i}, -1)

        grid[i][49] = "#"
        ast.FillTile(astar.Point{i, 49}, -1)
    }

    source := make([]astar.Point, 1)
    for {
        r := GetRandInt(50)
        c := GetRandInt(50)

        if grid[r][c] != "#" {
            grid[r][c] = "a"

            source[0].Row = r
            source[0].Col = c
            break
        }
    }

    var target astar.Point
    for {
        r := GetRandInt(50)
        c := GetRandInt(50)

        if grid[r][c] != "#" && grid[r][c] != "a" {
            grid[r][c] = "b"

            target.Row = r
            target.Col = c
            break
        }
    }

    PrintGrid(grid)
    end, _ := ast.FindPath(source, target, astar.RawDist)

    path := end
    for {
        if path.Row == source[0].Row && path.Col == source[0].Col {
            grid[path.Row][path.Col] = "A"
        } else if path.Row == target.Row && path.Col == target.Col {
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
    PrintGrid(grid)
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

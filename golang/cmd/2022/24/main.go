package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/ventsislav-georgiev/advent-of-code/golang/pkg/aoc"
)

var testInput = `#.######
#>>.<^<#
#.<..<<#
#>v.><>#
#<^v^^>#
######.#`

var testMode = false

func main() {
	aoc.Exec(task1, task2)
}

func task1(in io.Reader) {
	if testMode {
		in = strings.NewReader(testInput)
	}

	grid := parse(in)
	path := findPath(grid, Point{row: 0, col: 1}, Point{row: grid.rows, col: grid.cols - 1}, 0)

	println(path.dist)
	if testMode {
		path.print()
	}
}

func task2(in io.Reader) {
	if testMode {
		in = strings.NewReader(testInput)
	}

	grid := parse(in)
	src := Point{row: 0, col: 1}
	dest := Point{row: grid.rows, col: grid.cols - 1}

	path := findPath(grid, src, dest, 0)
	path = findPath(grid, dest, src, path.dist)
	path = findPath(grid, src, dest, path.dist)

	println(path.dist)
	if testMode {
		path.print()
	}
}

func findPath(grid *Grid, src Point, dest Point, dist int) *Path {
	queue := []*Path{}
	visited := map[Path]struct{}{}
	maxLen := (grid.rows - 1) * (grid.cols - 1)

	srcPath := &Path{pos: src, dist: dist}
	queue = append(queue, srcPath)
	visited[Path{pos: src, dist: srcPath.dist}] = struct{}{}

	var cur *Path
	for len(queue) > 0 {
		cur = queue[0]
		queue = queue[1:]

		if cur.pos == dest {
			break
		}

		adjcells := []Point{cur.pos, cur.pos.add(east), cur.pos.add(south), cur.pos.add(north), cur.pos.add(west)}
		dist := cur.dist + 1
		gridMap := grid.getMap(dist, maxLen)

		for _, p := range adjcells {
			wallCount := gridMap[p]
			if wallCount > 0 {
				continue
			}

			outbounds := p.row < 0 || p.row > grid.rows || p.col < 0 || p.col > grid.cols
			if outbounds {
				continue
			}

			if _, ok := visited[Path{pos: p, dist: dist}]; ok {
				continue
			}

			visited[Path{pos: p, dist: dist}] = struct{}{}

			path := &Path{
				pos:    p,
				dist:   dist,
				parent: cur,
			}

			queue = append(queue, path)
		}
	}

	return cur
}

var (
	north = Point{-1, 0}
	south = Point{1, 0}
	west  = Point{0, -1}
	east  = Point{0, 1}
)

type Grid struct {
	data map[Point]int
	rows int
	cols int
	bliz []Bliz
}

type Bliz struct {
	pos Point
	dir Point
}

type Point struct {
	row int
	col int
}

type Path struct {
	pos    Point
	dist   int
	parent *Path
}

func (p Point) add(dir Point) Point {
	return Point{row: p.row + dir.row, col: p.col + dir.col}
}

func (p *Path) print() {
	println(p.pos.row, p.pos.col)
	if p.parent != nil {
		p.parent.print()
	}
}

var mapCache = map[int]map[Point]int{}

func (g *Grid) getMap(steps int, maxSteps int) map[Point]int {
	steps = steps % maxSteps

	if gridMap, ok := mapCache[steps]; ok {
		return gridMap
	}

	for s := 1; s <= steps; s++ {
		if _, ok := mapCache[s]; ok {
			continue
		}

		for bIdx := 0; bIdx < len(g.bliz); bIdx++ {
			b := &g.bliz[bIdx]

			g.data[b.pos] -= 1
			if g.data[b.pos] < 0 {
				g.data[b.pos] = 0
			}

			b.pos = b.pos.add(b.dir)

			if b.pos.row == 0 {
				b.pos.row = g.rows - 1
			}
			if b.pos.row == g.rows {
				b.pos.row = 1
			}
			if b.pos.col == 0 {
				b.pos.col = g.cols - 1
			}
			if b.pos.col == g.cols {
				b.pos.col = 1
			}

			g.data[b.pos] += 1
		}

		gridMap := map[Point]int{}
		for k, v := range g.data {
			gridMap[k] = v
		}

		mapCache[steps] = gridMap
	}

	if testMode {
		println(steps, "steps")
		printGrid(mapCache[steps], g.rows, g.cols)
	}

	return mapCache[steps]
}

func parse(in io.Reader) *Grid {
	grid := &Grid{
		data: make(map[Point]int),
		bliz: make([]Bliz, 0),
		rows: 0,
		cols: 0,
	}

	scanner := bufio.NewScanner(in)
	var row int
	for scanner.Scan() {
		line := scanner.Bytes()

		if grid.rows == 0 {
			grid.cols = len(line) - 1
		}

		for col, ch := range line {
			pKey := Point{row, col}
			grid.data[pKey] = 0

			switch ch {
			case '#':
				grid.data[pKey] += 1
			case '>', '<', '^', 'v':
				grid.data[pKey] += 1

				var dir Point
				switch ch {
				case '>':
					dir = east
				case '<':
					dir = west
				case '^':
					dir = north
				case 'v':
					dir = south
				}

				bliz := Bliz{Point{row, col}, dir}
				grid.bliz = append(grid.bliz, bliz)
			}
		}
		row++
	}

	grid.rows = row - 1
	return grid
}

func printGrid(blizMap map[Point]int, rows, cols int) {
	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {
			fill, ok := blizMap[Point{row, col}]
			if ok {
				print(fill)
			} else {
				print(".")
			}
		}
		println()
	}
	println()
}

package engine

import (
	"math/rand"
	"sort"
	"time"
)

const (
	rowsCnt = 4
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type value int64

type valuable interface {
	NumValue() value
}

type valuableReserve struct {
	Board *Board
	Dist  int
}

func (v *valuableReserve) NumValue() value {
	return value(v.Board.LineLength - v.Dist)
}

type valuableEmpty struct {
	Board *Board
}

func (v *valuableEmpty) NumValue() value {
	return value(v.Board.LineLength * 2 /*sides*/ * rowsCnt)
}

type valuableOwn struct {
	Board *Board
	Which int
}

func (v *valuableOwn) NumValue() value {
	tmp := valuableEmpty{Board: v.Board}
	emptyVal := tmp.NumValue()
	var ret value
	for w := value(1); w <= value(v.Which); w++ {
		factor := value(1)
		for i := value(0); i < w; i++ {
			factor *= rowsCnt * (value(v.Board.LineLength) - i)
		}
		ret += value(emptyVal * factor)
	}
	return ret
}

func (b *Board) evaluateLine(rowProvider *rowProvider, index int, ownType PointType) []valuable {
	ret := []valuable{}
	ownCnt := 0
	for i := index; i < (index + b.LineLength); i++ {
		cellAff := getCellAff(rowProvider.getCellType(i), ownType)
		var v valuable
		switch cellAff {
		case cellAffEmpty:
			v = &valuableEmpty{
				Board: b,
			}
		case cellAffOwn:
			ownCnt++
			v = &valuableOwn{
				Board: b,
				Which: ownCnt,
			}
		case cellAffEnemy:
			return []valuable{}
		}
		ret = append(ret, v)
	}

	//TODO: calc reserve

	return ret
}

func invertPointType(pt PointType) PointType {
	switch pt {
	case PointX:
		return PointO
	case PointO:
		return PointX
	default:
		panic("wrong point type")
	}
}

func (b *Board) evaluatePos(x, y int, ownType PointType) (ownVals []valuable, enemyVals []valuable) {
	if b.Cells[y][x] != CellEmpty {
		panic("evaluatePos is called for non-empty cell")
	}
	// temporary put own point into the given cell
	defer func() { b.Cells[y][x] = CellEmpty }()

	rowProviders := []*rowProvider{
		makeRowProviderHor(b, x, y),
		makeRowProviderSlash(b, x, y),
		makeRowProviderVert(b, x, y),
		makeRowProviderBackslash(b, x, y),
	}

	ownVals = []valuable{}
	enemyVals = []valuable{}

	for _, rp := range rowProviders {
		for i := -b.LineLength + 1; i <= 0; i++ {
			b.Cells[y][x] = getCellTypeByAff(cellAffOwn, ownType)
			ownVals = append(ownVals, b.evaluateLine(rp, i, ownType)...)
			b.Cells[y][x] = getCellTypeByAff(cellAffEnemy, ownType)
			enemyVals = append(enemyVals, b.evaluateLine(rp, i, invertPointType(ownType))...)
		}
	}

	return ownVals, enemyVals
}

func valuablesToNum(valuables []valuable, factor float64) value {
	var ret value
	for _, v := range valuables {
		ret += v.NumValue()
	}
	if factor != 1 {
		ret = value(float64(ret) * factor)
	}
	return ret
}

func allValuablesToNum(ownVals []valuable, enemyVals []valuable) value {
	// We consider both ownVals and enemyVals, but we apply a 0.9 factor to
	// enemyVals to make it slightly more important to build our own stuff instead
	// of breaking enemy's stuff
	return valuablesToNum(ownVals, 1) + valuablesToNum(enemyVals, 0.9)
}

type pos struct {
	X, Y      int
	Value     value
	OwnVals   []valuable
	EnemyVals []valuable
}

type positions []pos

func (positions positions) Len() int {
	return len(positions)
}

func (positions positions) Less(i, j int) bool {
	return positions[i].Value < positions[j].Value
}

func (positions positions) Swap(i, j int) {
	positions[i], positions[j] = positions[j], positions[i]
}

func (b *Board) EvaluateAllPos(ownType PointType) positions {
	positions := positions{}

	for y := range b.Cells {
		for x, c := range b.Cells[y] {
			if c == CellEmpty {
				curPos := pos{
					X: x,
					Y: y,
				}

				ownVals, enemyVals := b.evaluatePos(x, y, ownType)

				curPos.OwnVals = ownVals
				curPos.EnemyVals = enemyVals
				curPos.Value = allValuablesToNum(ownVals, enemyVals)

				positions = append(positions, curPos)
			}
		}
	}

	for i := range positions {
		j := rand.Intn(i + 1)
		positions[i], positions[j] = positions[j], positions[i]
	}

	sort.Sort(positions)

	return positions
}

func (b *Board) FindBestPos(ownType PointType) *pos {
	var ret *pos
	positions := b.EvaluateAllPos(ownType)
	if len(positions) > 0 {
		ret = &positions[len(positions)-1]
	}
	return ret
}

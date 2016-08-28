package engine

type CellType int

const (
	CellEmpty CellType = iota
	CellX
	CellO
	CellXStriked
	CellOStriked
	CellWall
)

type PointType int

const (
	PointX PointType = iota
	PointO
)

type Board struct {
	Cells      [][]CellType
	LineLength int
}

func MakeBoard(sizeX, sizeY int) *Board {
	board := Board{
		LineLength: 5,
	}
	board.Cells = make([][]CellType, sizeY)
	for y := range board.Cells {
		board.Cells[y] = make([]CellType, sizeX)
		for x := range board.Cells[y] {
			board.Cells[y][x] = CellEmpty
		}
	}
	return &board
}

func (b *Board) Put(x, y int, ownType PointType) {
	cellType := getCellTypeByAff(cellAffOwn, ownType)
	var cellStrikedType CellType
	switch cellType {
	case CellX:
		cellStrikedType = CellXStriked
	case CellO:
		cellStrikedType = CellOStriked
	default:
		panic("")
	}

	b.Cells[y][x] = cellType

	// check if it's win
	// TODO: refactor
	rowProviders := []*rowProvider{
		makeRowProviderHor(b, x, y),
		makeRowProviderSlash(b, x, y),
		makeRowProviderVert(b, x, y),
		makeRowProviderBackslash(b, x, y),
	}

	for _, rp := range rowProviders {
		var a1, a2 int
		for i := 0; i > -b.LineLength; i-- {
			if rp.getCellType(i) == cellType {
				a1++
			} else {
				break
			}
		}

		for i := 1; i < b.LineLength; i++ {
			if rp.getCellType(i) == cellType {
				a2++
			} else {
				break
			}
		}

		if a1+a2 == b.LineLength {
			//win!
			for i := 0; i < a1; i++ {
				rp.setCellType(-i, cellStrikedType)
			}
			for i := 1; i <= a2; i++ {
				rp.setCellType(i, cellStrikedType)
			}
		} else if a1+a2 > b.LineLength {
			panic("a1 + a2 should not be > b.LineLength")
		}
	}
}

type cellAff int

const (
	cellAffEmpty = iota
	cellAffOwn
	cellAffEnemy //TODO: rename, because it actually includes walls, etc
)

func getCellAff(cellType CellType, ownType PointType) cellAff {
	switch cellType {
	case CellEmpty:
		return cellAffEmpty
	case CellX:
		if ownType == PointX {
			return cellAffOwn
		} else {
			return cellAffEnemy
		}
	case CellO:
		if ownType == PointO {
			return cellAffOwn
		} else {
			return cellAffEnemy
		}
	case CellXStriked, CellOStriked, CellWall:
		return cellAffEnemy
	default:
		panic("unknown cellType")
	}
}

func getCellTypeByAff(cellAff cellAff, ownType PointType) CellType {
	switch cellAff {
	case cellAffEmpty:
		return CellEmpty
	case cellAffOwn:
		if ownType == PointX {
			return CellX
		} else {
			return CellO
		}
	case cellAffEnemy:
		if ownType == PointX {
			return CellO
		} else {
			return CellX
		}
	default:
		panic("unknown cellAff")
	}
}

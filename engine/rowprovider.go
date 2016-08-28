package engine

type rowProvider struct {
	Board   *Board
	CenterX int
	CenterY int
	XStep   int
	YStep   int
}

func (rp *rowProvider) getCellType(index int) CellType {
	x := rp.CenterX + index*rp.XStep
	y := rp.CenterY + index*rp.YStep
	if y >= 0 && y < len(rp.Board.Cells) && x >= 0 && x < len(rp.Board.Cells[y]) {
		return rp.Board.Cells[y][x]
	} else {
		return CellWall
	}
}

func (rp *rowProvider) setCellType(index int, cellType CellType) {
	x := rp.CenterX + index*rp.XStep
	y := rp.CenterY + index*rp.YStep
	if y >= 0 && y < len(rp.Board.Cells) && x >= 0 && x < len(rp.Board.Cells[y]) {
		rp.Board.Cells[y][x] = cellType
	} else {
		panic("can't set cell type of wall")
	}
}

func makeRowProviderHor(board *Board, centerX, centerY int) *rowProvider {
	return &rowProvider{
		Board:   board,
		CenterX: centerX,
		CenterY: centerY,
		XStep:   1,
		YStep:   0,
	}
}

func makeRowProviderSlash(board *Board, centerX, centerY int) *rowProvider {
	return &rowProvider{
		Board:   board,
		CenterX: centerX,
		CenterY: centerY,
		XStep:   1,
		YStep:   -1,
	}
}

func makeRowProviderVert(board *Board, centerX, centerY int) *rowProvider {
	return &rowProvider{
		Board:   board,
		CenterX: centerX,
		CenterY: centerY,
		XStep:   0,
		YStep:   -1,
	}
}
func makeRowProviderBackslash(board *Board, centerX, centerY int) *rowProvider {
	return &rowProvider{
		Board:   board,
		CenterX: centerX,
		CenterY: centerY,
		XStep:   -1,
		YStep:   -1,
	}
}

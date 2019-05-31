package internal

type errGameNoBoard struct {
	message string
}

func newErrGameNoBoard() *errGameNoBoard {
	return &errGameNoBoard{"no board"}
}

func (e *errGameNoBoard) Error() string {
	return e.message
}

type errGameNoFigureInCell struct {
	message string
}

func newErrGameNoFigureInCell() *errGameNoFigureInCell {
	return &errGameNoFigureInCell{"no figure in the cell"}
}

func (e *errGameNoFigureInCell) Error() string {
	return e.message
}

type errGameNotFigureOwner struct {
	message string
}

func newErrGameNotFigureOwner() *errGameNotFigureOwner {
	return &errGameNotFigureOwner{"you are not UnitOwner of this figure"}
}

func (e *errGameNotFigureOwner) Error() string {
	return e.message
}

type errGameFigureNotMovable struct {
	message string
}

func newErrGameFigureNotMovable() *errGameFigureNotMovable {
	return &errGameFigureNotMovable{"figure already set"}
}

func (e *errGameFigureNotMovable) Error() string {
	return e.message
}

type errFigureMoveOutOfStartZone struct {
	message string
}

func newErrFigureMoveOutOfStartZone() *errFigureMoveOutOfStartZone {
	return &errFigureMoveOutOfStartZone{"figure out of start zone"}
}

func (e *errFigureMoveOutOfStartZone) Error() string {
	return e.message
}

type errCellNotEmpty struct {
	message string
}

func newErrCellNotEmpty() *errCellNotEmpty {
	return &errCellNotEmpty{"cell not empty"}
}

func (e *errCellNotEmpty) Error() string {
	return e.message
}

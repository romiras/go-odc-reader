// Package visitor defines the Visitor pattern interface for traversing document trees.
package visitor

// Visitor defines the interface for traversing and processing document elements.
// This is the Visitor role in the Visitor design pattern.
type Visitor interface {
	// PartStart signals the beginning of a "part" (container element).
	// This can be an StdTextModel, ViewPiece, or another container.
	PartStart()

	// PartEnd signals the end of a "part".
	PartEnd()

	// FoldLeft signals a left fold marker has been found.
	// If collapsed is true, the first part that follows is the "hidden" part,
	// otherwise the first part is the "alternative" text.
	FoldLeft(collapsed bool)

	// FoldRight signals a right fold marker has been found.
	FoldRight()

	// TextShortPiece processes a text piece with 8-bit characters (Latin-1).
	// The piece parameter will be a *textmodel.ShortPiece that can be type-asserted.
	TextShortPiece(piece interface{})

	// TextLongPiece processes a text piece with 16-bit characters (Unicode).
	// The piece parameter will be a *textmodel.LongPiece that can be type-asserted.
	TextLongPiece(piece interface{})
}

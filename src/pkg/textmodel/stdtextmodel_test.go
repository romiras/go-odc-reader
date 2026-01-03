package textmodel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

// MockReader implements store.Reader for testing
type MockReader struct {
	data   []byte
	pos    int
	stores []store.Store
}

func NewMockReader(data []byte) *MockReader {
	return &MockReader{data: data, stores: make([]store.Store, 0)}
}

func (m *MockReader) ReadByte() (byte, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	b := m.data[m.pos]
	m.pos++
	return b, nil
}

func (m *MockReader) ReadSignedByte() (oberon.Byte, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	b := m.data[m.pos]
	m.pos++
	return oberon.Byte(b), nil
}

func (m *MockReader) ReadInt() (oberon.Integer, error) {
	if m.pos+4 > len(m.data) {
		return 0, io.EOF
	}
	var val int32
	buf := bytes.NewReader(m.data[m.pos : m.pos+4])
	binary.Read(buf, binary.LittleEndian, &val)
	m.pos += 4
	return oberon.Integer(val), nil
}

func (m *MockReader) ReadSInt() (oberon.ShortInt, error) {
	return 0, fmt.Errorf("not implemented")
}

func (m *MockReader) ReadSChar() (oberon.ShortChar, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	b := m.data[m.pos]
	m.pos++
	return oberon.ShortChar(b), nil
}

func (m *MockReader) ReadLChar() (oberon.Char, error) {
	if m.pos+2 > len(m.data) {
		return 0, io.EOF
	}
	var val uint16
	buf := bytes.NewReader(m.data[m.pos : m.pos+2])
	binary.Read(buf, binary.LittleEndian, &val)
	m.pos += 2
	return oberon.Char(val), nil
}

func (m *MockReader) ReadSString() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (m *MockReader) ReadVersion(min, max oberon.Integer) (oberon.Integer, error) {
	b, err := m.ReadSignedByte()
	if err != nil {
		return 0, err
	}
	v := oberon.Integer(b)
	if v < min || v > max {
		return v, fmt.Errorf("version mismatch")
	}
	return v, nil
}

// MockStore implements store.Store query methods
type MockStore struct {
	store.BaseStore
}

func (ms *MockStore) GetTypeName() string              { return "MockStore" }
func (ms *MockStore) GetTypePath() store.TypePath      { return store.TypePath{"MockStore"} }
func (ms *MockStore) Internalize(r store.Reader) error { return nil }

func (m *MockReader) ReadStore() (store.Store, error) {
	// For testing attributes, return a dummy store
	return &MockStore{BaseStore: store.NewBaseStore(0)}, nil
}

func (m *MockReader) IsCancelled() bool {
	return false
}

func (m *MockReader) TurnIntoAlien(cause int) error {
	return fmt.Errorf("turned into alien")
}

func TestStdTextModel_Internalize_Empty(t *testing.T) {
	// Hierarchy: StdTextModel -> TextModel -> ContainerModel -> Model -> Elem -> BaseStore
	// Each level reads a version byte. Total 6 levels.
	// Versions: 0, 0, 0, 0, 0, 0
	// MetaLen: 0 (4 bytes)
	// Ano: -1 (0xFF)
	data := []byte{
		0, 0, 0, 0, 0, 0, // Versions
		0, 0, 0, 0, // MetaLen
		0xFF, // Ano (End)
	}
	reader := NewMockReader(data)

	model := NewStdTextModel(0)
	err := model.Internalize(reader)

	if err != nil {
		t.Fatalf("Internalize failed: %v", err)
	}

	if len(model.pieces) != 0 {
		t.Errorf("Expected 0 pieces, got %d", len(model.pieces))
	}
}

func TestStdTextModel_Internalize_ShortPiece(t *testing.T) {
	// Construct data
	buf := new(bytes.Buffer)

	// Versions (6 levels)
	for i := 0; i < 6; i++ {
		buf.WriteByte(0)
	}

	// MetaLen: 4 bytes (0)
	binary.Write(buf, binary.LittleEndian, int32(0))

	// Ano: 0 (new attribute)
	buf.WriteByte(0)

	// Attribute store: MockReader.ReadStore consumes NO bytes

	// PieceLen: 5 (Short piece "Hello")
	binary.Write(buf, binary.LittleEndian, int32(5))

	// Next Ano: -1 (EndOfPieces)
	buf.WriteByte(0xFF)

	// Piece Content: "Hello"
	buf.WriteString("Hello")

	reader := NewMockReader(buf.Bytes())
	model := NewStdTextModel(1)

	err := model.Internalize(reader)
	if err != nil {
		t.Fatalf("Internalize failed: %v", err)
	}

	if len(model.pieces) != 1 {
		t.Fatalf("Expected 1 piece, got %d", len(model.pieces))
	}

	sp, ok := model.pieces[0].(*ShortPiece)
	if !ok {
		t.Fatalf("Expected ShortPiece")
	}

	content := string(sp.buffer[:sp.length])
	if content != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", content)
	}
}

// Package reader provides binary file reading for .odc documents.
package reader

import (
	"encoding/binary"
	"fmt"
	"io"

	"odcread/pkg/alien"
	"odcread/pkg/oberon"
	"odcread/pkg/store"
	"odcread/pkg/typeregister"
)

// Cause codes for alien conversion
const (
	TypeNotFound = 1 // Type not registered
	AlienVersion = 2 // Version out of range
)

// TypeEntry represents a type in the type dictionary.
type TypeEntry struct {
	Name   string
	BaseID oberon.Integer
}

// ReaderState stores the reader's position state.
type ReaderState struct {
	Next int64 // Position of next store
	End  int64 // Position after last store
}

// Reader reads binary .odc format and manages parsing state.
type Reader struct {
	rider        io.ReadSeeker
	cancelled    bool
	cause        int
	readAlien    bool
	typeList     []*TypeEntry
	elemList     []store.Store
	storeList    []store.Store
	currentStore store.Store
	state        *ReaderState
}

// NewReader creates a new Reader for the given input stream.
func NewReader(r io.ReadSeeker) *Reader {
	return &Reader{
		rider:     r,
		typeList:  make([]*TypeEntry, 0),
		elemList:  make([]store.Store, 0),
		storeList: make([]store.Store, 0),
		state:     &ReaderState{},
	}
}

// ReadSChar reads a single 8-bit character.
func (r *Reader) ReadSChar() (oberon.ShortChar, error) {
	var ch oberon.ShortChar
	err := binary.Read(r.rider, binary.LittleEndian, &ch)
	return ch, err
}

// ReadLChar reads a single 16-bit character.
func (r *Reader) ReadLChar() (oberon.Char, error) {
	var ch oberon.Char
	err := binary.Read(r.rider, binary.LittleEndian, &ch)
	return ch, err
}

// ReadByte reads a single unsigned byte (implements io.ByteReader).
func (r *Reader) ReadByte() (byte, error) {
	var b byte
	err := binary.Read(r.rider, binary.LittleEndian, &b)
	return b, err
}

// ReadSignedByte reads a single signed byte.
func (r *Reader) ReadSignedByte() (oberon.Byte, error) {
	var b oberon.Byte
	err := binary.Read(r.rider, binary.LittleEndian, &b)
	return b, err
}

// ReadSInt reads a 16-bit signed integer.
func (r *Reader) ReadSInt() (oberon.ShortInt, error) {
	var val oberon.ShortInt
	err := binary.Read(r.rider, binary.LittleEndian, &val)
	return val, err
}

// ReadInt reads a 32-bit signed integer.
func (r *Reader) ReadInt() (oberon.Integer, error) {
	var val oberon.Integer
	err := binary.Read(r.rider, binary.LittleEndian, &val)
	return val, err
}

// ReadSString reads a null-terminated short string.
func (r *Reader) ReadSString() (string, error) {
	var chars []oberon.ShortChar
	for {
		ch, err := r.ReadSChar()
		if err != nil {
			return "", err
		}
		if ch == 0 {
			break
		}
		chars = append(chars, ch)
	}
	return string(chars), nil
}

// ReadVersion reads and validates a version byte.
// If the version is not in [min, max], the current store is turned into an alien.
func (r *Reader) ReadVersion(min, max oberon.Integer) (oberon.Integer, error) {
	versionByte, err := r.ReadSignedByte()
	if err != nil {
		return 0, err
	}

	version := oberon.Integer(versionByte)

	if version < min || version > max {
		r.TurnIntoAlien(AlienVersion)
		return version, fmt.Errorf("version %d out of range [%d, %d]", version, min, max)
	}

	return version, nil
}

// IsCancelled returns whether the current read has been cancelled.
func (r *Reader) IsCancelled() bool {
	return r.cancelled
}

// TurnIntoAlien cancels the current read and marks it for alien conversion.
func (r *Reader) TurnIntoAlien(cause int) error {
	r.cancelled = true
	r.cause = cause
	r.readAlien = true
	return fmt.Errorf("turned into alien: cause %d", cause)
}

// ReadStore reads a store from the binary stream.
func (r *Reader) ReadStore() (store.Store, error) {
	return r.readStoreOrElemStore()
}

// readStoreOrElemStore reads either a Store or Elem-type store.
func (r *Reader) readStoreOrElemStore() (store.Store, error) {
	// Read the store marker
	marker, err := r.ReadSChar()
	if err != nil {
		return nil, fmt.Errorf("failed to read store marker: %w", err)
	}

	switch marker {
	case store.NIL:
		return r.readNilStore()
	case store.LINK:
		return r.readLinkStore()
	case store.NEWLINK:
		return r.readNewLinkStore()
	case store.STORE, store.ELEM:
		return r.readNewStore(marker == store.ELEM)
	default:
		return nil, fmt.Errorf("unknown store marker: 0x%X", marker)
	}
}

// readNilStore handles nil store markers.
func (r *Reader) readNilStore() (store.Store, error) {
	// Nil stores still have header fields that must be consumed
	comment, err := r.ReadInt() // 4 bytes
	if err != nil {
		return nil, fmt.Errorf("failed to read comment: %w", err)
	}

	next, err := r.ReadInt() // 4 bytes
	if err != nil {
		return nil, fmt.Errorf("failed to read next: %w", err)
	}

	// Update state tracking
	currentPos, err := r.rider.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	r.state.End = currentPos

	// Calculate next pointer
	if next > 0 || (next == 0 && comment%2 == 1) {
		r.state.Next = r.state.End + int64(next)
	} else {
		r.state.Next = 0
	}

	return nil, nil
}

// readLinkStore reads a link to an Elem-type store.
func (r *Reader) readLinkStore() (store.Store, error) {
	// LINK stores have full headers: id, comment, next (12 bytes total)
	// From Component Pascal: rd.ReadInt(id); rd.ReadInt(comment); rd.ReadInt(next);
	id, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read link ID: %w", err)
	}

	comment, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read comment: %w", err)
	}

	next, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read next: %w", err)
	}

	// Update state tracking (same logic as NIL stores)
	currentPos, err := r.rider.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	r.state.End = currentPos

	// Calculate next pointer
	if next > 0 || (next == 0 && comment%2 == 1) {
		r.state.Next = r.state.End + int64(next)
	} else {
		r.state.Next = 0
	}

	// Look up in elem list
	if id < 0 || int(id) >= len(r.elemList) {
		return nil, fmt.Errorf("invalid elem link ID: %d", id)
	}

	return r.elemList[id], nil
}

// readNewLinkStore reads a link to a non-Elem-type store.
func (r *Reader) readNewLinkStore() (store.Store, error) {
	// NEWLINK stores have full headers: id, comment, next (12 bytes total)
	// From Component Pascal: rd.ReadInt(id); rd.ReadInt(comment); rd.ReadInt(next);
	id, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read new link ID: %w", err)
	}

	comment, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read comment: %w", err)
	}

	next, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read next: %w", err)
	}

	// Update state tracking (same logic as NIL stores)
	currentPos, err := r.rider.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	r.state.End = currentPos

	// Calculate next pointer
	if next > 0 || (next == 0 && comment%2 == 1) {
		r.state.Next = r.state.End + int64(next)
	} else {
		r.state.Next = 0
	}

	// Look up in store list
	if id < 0 || int(id) >= len(r.storeList) {
		return nil, fmt.Errorf("invalid store link ID: %d", id)
	}

	return r.storeList[id], nil
}

// readNewStore reads a new store (not a link).
func (r *Reader) readNewStore(isElem bool) (store.Store, error) {
	// Calculate the store ID
	id := oberon.Integer(len(r.elemList))
	if !isElem {
		id = oberon.Integer(len(r.storeList))
	}

	// Read the type path
	path, err := r.readPath()
	if err != nil {
		return nil, fmt.Errorf("failed to read type path: %w", err)
	}

	if len(path) == 0 {
		return nil, fmt.Errorf("empty type path")
	}

	// Get the type name (first element of path - root type)
	typeName := path[0]

	// Read the store header fields
	_, err = r.ReadInt() // comment (not used)
	if err != nil {
		return nil, fmt.Errorf("failed to read comment: %w", err)
	}

	pos1, err := r.rider.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	next, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read next: %w", err)
	}

	down, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read down: %w", err)
	}

	length, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	pos, err := r.rider.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get position after header: %w", err)
	}

	// Calculate state positions
	if next > 0 {
		r.state.Next = pos1 + int64(next) + 4
	} else {
		r.state.Next = 0
	}

	var downPos int64
	if down > 0 {
		downPos = pos1 + int64(down) + 8
	}

	r.state.End = pos + int64(length)
	r.cause = 0

	// Try to create a store instance from the type registry
	proxy := typeregister.GetInstance().Get(typeName)
	var st store.Store

	if proxy != nil {
		st = proxy.NewInstance(id)
	} else {
		r.cause = TypeNotFound
	}

	// If we successfully created a store, try to internalize it
	if st != nil {
		// Save the current store's end position BEFORE swapping states
		storeEnd := r.state.End

		// Save the current state and create new state for nested reads
		saveState := r.state
		r.state = &ReaderState{}

		// Internalize the store
		st.Internalize(r)

		// Restore the state
		r.state = saveState

		// If internalization failed, turn it into an alien
		if r.cause != 0 {
			st = nil
		} else {
			// Verify we're at the expected position using the SAVED end position
			currentPos, _ := r.rider.Seek(0, io.SeekCurrent)
			if currentPos != storeEnd {
				return nil, fmt.Errorf("position mismatch after internalize: expected %d, got %d", storeEnd, currentPos)
			}
		}
	}

	// If we have a valid store, add it to the appropriate list
	if st != nil {
		if isElem {
			r.elemList = append(r.elemList, st)
		} else {
			r.storeList = append(r.storeList, st)
		}
		return st, nil
	}

	// If we failed to create or internalize the store, create an alien
	r.rider.Seek(pos, io.SeekStart)

	alienStore := alien.NewAlien(id, path)

	if isElem {
		r.elemList = append(r.elemList, alienStore)
	} else {
		r.storeList = append(r.storeList, alienStore)
	}

	// Save the store's end position BEFORE swapping states
	// This is critical for nested aliens - we need the actual end position, not the empty state's End
	storeEnd := r.state.End

	// Save state and internalize the alien
	saveState := r.state
	r.state = &ReaderState{}

	err = r.internalizeAlien(alienStore, downPos, storeEnd)
	if err != nil {
		r.state = saveState
		return nil, fmt.Errorf("failed to internalize alien: %w", err)
	}

	r.state = saveState

	// Verify position after alien internalization using the SAVED end position
	currentPos, _ := r.rider.Seek(0, io.SeekCurrent)
	if currentPos != storeEnd {
		return nil, fmt.Errorf("position mismatch after alien: expected %d, got %d", storeEnd, currentPos)
	}

	// Reset state after reading alien
	r.cause = 0
	r.cancelled = false
	r.readAlien = true

	return alienStore, nil
}

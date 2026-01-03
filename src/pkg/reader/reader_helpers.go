package reader

import (
	"fmt"
	"io"
	"strings"

	"odcread/pkg/alien"
	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

// fixTypeName replaces "Desc" suffix with "^" to match Oberon naming conventions.
func fixTypeName(name string) string {
	if len(name) >= 4 && strings.HasSuffix(name, "Desc") {
		return name[:len(name)-4] + "^"
	}
	return name
}

// addPathComponent adds a type name to the type dictionary.
func (r *Reader) addPathComponent(first bool, typeName string) {
	next := len(r.typeList)
	curr := next - 1
	if !first {
		r.typeList[curr].BaseID = oberon.Integer(next)
	}
	r.typeList = append(r.typeList, &TypeEntry{
		Name:   typeName,
		BaseID: -1,
	})
}

// readPath reads the type path from the binary stream.
// This is the CORRECTED implementation matching the C++ logic.
func (r *Reader) readPath() (store.TypePath, error) {
	var path store.TypePath

	// Read the first marker
	marker, err := r.ReadSChar()
	if err != nil {
		return nil, fmt.Errorf("failed to read path marker: %w", err)
	}

	// Loop through NEWEXT markers (this was the bug - it wasn't looping!)
	i := 0
	for marker == store.NEWEXT {
		// Read the type name string
		typeName, err := r.ReadSString()
		if err != nil {
			return nil, fmt.Errorf("failed to read type name: %w", err)
		}

		typeName = fixTypeName(typeName)
		path = append(path, typeName)
		r.addPathComponent(i == 0, typeName)
		i++

		// Read the next marker (this is critical - was missing in buggy version!)
		marker, err = r.ReadSChar()
		if err != nil {
			return nil, fmt.Errorf("failed to read next marker: %w", err)
		}
	}

	if marker == store.NEWBASE {
		// Read the base type name
		typeName, err := r.ReadSString()
		if err != nil {
			return nil, fmt.Errorf("failed to read base type name: %w", err)
		}

		typeName = fixTypeName(typeName)
		path = append(path, typeName)
		r.addPathComponent(i == 0, typeName)

		return path, nil

	} else if marker == store.OLDTYPE {
		// Read the type ID and traverse the type dictionary chain
		typeID, err := r.ReadInt()
		if err != nil {
			return nil, fmt.Errorf("failed to read type ID: %w", err)
		}

		// Update the base ID if we have previous entries
		if i > 0 {
			r.typeList[len(r.typeList)-1].BaseID = typeID
		}

		// Loop through the entire type dictionary chain until baseId == -1
		// (This was also a bug - it only read ONE type instead of looping!)
		for typeID != -1 {
			if typeID < 0 || int(typeID) >= len(r.typeList) {
				return nil, fmt.Errorf("invalid type ID: %d", typeID)
			}

			path = append(path, r.typeList[typeID].Name)
			typeID = r.typeList[typeID].BaseID
			i++
		}

		return path, nil
	}

	return nil, fmt.Errorf("unexpected path marker: 0x%X", marker)
}

// internalizeAlien reads the contents of an alien store.
func (r *Reader) internalizeAlien(alienStore *alien.Alien, down, end int64) error {
	next := down
	if next == 0 {
		next = end
	}

	for {
		currentPos, err := r.rider.Seek(0, io.SeekCurrent)
		if err != nil {
			return fmt.Errorf("failed to get current position: %w", err)
		}

		if currentPos >= end {
			break
		}

		if currentPos < next {
			// Read a piece (unstructured binary data)
			length := next - currentPos
			buf := make([]byte, length)
			n, err := r.rider.Read(buf)
			if err != nil {
				return fmt.Errorf("failed to read alien piece: %w", err)
			}
			if int64(n) != length {
				return fmt.Errorf("short read: expected %d bytes, got %d", length, n)
			}

			piece := alien.NewAlienPiece(buf)
			alienStore.AddComponent(piece)

		} else {
			// Seek to the store position and read it
			_, err := r.rider.Seek(next, io.SeekStart)
			if err != nil {
				return fmt.Errorf("failed to seek to store: %w", err)
			}

			st, err := r.ReadStore()
			if err != nil {
				return fmt.Errorf("failed to read store in alien: %w", err)
			}

			part := alien.NewAlienPart(st)
			alienStore.AddComponent(part)

			// Update next position
			if r.state.Next > 0 {
				next = r.state.Next
			} else {
				next = end
			}
		}
	}

	return nil
}

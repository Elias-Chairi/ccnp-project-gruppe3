package util

import "net"

// Generates a unique ID based on the current number of existing IDs.
// Is a id generator that maps the created ID to the net.Conn to be implemented by node and control panel entities.
func GetUniqueID(existingIDs map[uint8]net.Conn) uint8 {
	newID := uint8(len(existingIDs) + 1)
	for {
		if _, exists := existingIDs[newID]; !exists {
			break
		}
		newID++
	}
	return newID
}

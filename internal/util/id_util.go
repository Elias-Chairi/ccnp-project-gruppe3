package util

import "net"

// GetUniqueID generates a unique uint8 ID not present in the existingIDs map.
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

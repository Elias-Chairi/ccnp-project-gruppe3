package general

type number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func GetUniqueMapKey[T any, V number](m map[V]T) V {
	newKey := V(1)
	for {
		if _, exists := m[newKey]; !exists {
			break
		}
		newKey++
	}
	return newKey
}

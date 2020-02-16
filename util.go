package caching

import "hash/fnv"

func generateKey(data string) string {
	hash := fnv.New64a()
	hash.Write([]byte(data))

	return string(hash.Sum64())
}

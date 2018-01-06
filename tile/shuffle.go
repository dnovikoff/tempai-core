package tile

import "math/rand"

func Shuffle(tiles Instances, rnd *rand.Rand) {
	n := len(tiles)
	for i := n - 1; i > 0; i-- {
		j := rnd.Intn(i)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
	return
}

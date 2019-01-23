package plugin

// Represents chunk indexes <first, last)
type Range struct {
	First int
	Last  int
}

// Calculate chunk indexes when dividing large slices
// Example: for length=50, chunkLength=20, result will be = {0, 20},{20,40},{40,50}
func CalculateChunkIndexes(length int, chunkLength int) []Range {
	var indexes []Range

	for i := 0; i < length; i += chunkLength {
		minIndex := i
		maxIndex := i + chunkLength

		if maxIndex > length {
			maxIndex = length
		}

		indexes = append(indexes, Range{minIndex, maxIndex})
	}

	return indexes
}

// Chunk metrics into smaller parts with defined maximum length
func ChunkMetrics(mts []Metric, chunkLength int) [][]Metric {
	var result [][]Metric

	for _, indexes := range CalculateChunkIndexes(len(mts), chunkLength) {
		result = append(result, mts[indexes.First:indexes.Last])
	}

	return result
}

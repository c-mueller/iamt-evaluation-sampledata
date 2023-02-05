package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	numRows           = flag.Int("num-rows", 10000, "The number of rows datasets to generate")
	dimensionCount    = flag.Int("num-dims", 2, "The number of input dimensions (e.g. sensors)")
	outlierPercentage = flag.Float64("outlier-percentage", 0.1, "The percentage (in decimal) of outliers in the output dataset")

	minClusterCount = flag.Int("min-clusters", 1, "Minimum number of value clusters per dimension")
	maxClusterCount = flag.Int("max-clusters", 10, "Maximum number of value clusters per dimension")

	clusterMax    = flag.Float64("max-value", 300, "Maximum value to use for inliers")
	clusterSpread = flag.Float64("cluster-spread", 5, "The spread of values in one cluster")

	generateNegativeOutliers = flag.Bool("negative-outliers", false, "Generate negative outliers")

	shuffle = flag.Bool("shuffle", false, "Shuffle output randomly")

	humanReadableFormat = flag.Bool("human", false, "Print in human Readable format")

	multiDimensionalOutliersOnly  = flag.Bool("only-multidimensional-outliers", false, "Only generate multidimenstional outliers")
	singleDimensionalOutliersOnly = flag.Bool("only-singledimenstional-outliers", false, "Only generate single dimensional Outliers")

	seed = flag.Int64("seed", time.Now().UnixNano(), "Random Seed")
)

func init() {
	flag.Parse()
}

type dataRow struct {
	IsOutlier    bool
	Measurements []float64
}

func (r dataRow) printHeader() {
	line := ""
	for i, _ := range r.Measurements {
		line = fmt.Sprintf("%s%15s", line, fmt.Sprintf("DIM-%d", i+1))
	}
	line = fmt.Sprintf("%s %15s", line, "IS OUTLIER")
	fmt.Println(line)
}

func (r dataRow) printCsvHeader() {
	line := ""
	for i, _ := range r.Measurements {
		line = fmt.Sprintf("%s%s,", line, fmt.Sprintf("dim-%d", i+1))
	}
	line = fmt.Sprintf("%s%s", line, "outlier")
	fmt.Println(line)
}

func (r dataRow) printCsvRow() {
	line := ""
	for _, measurement := range r.Measurements {
		line = fmt.Sprintf("%s%f,", line, measurement)
	}
	line = fmt.Sprintf("%s%v", line, r.IsOutlier)
	fmt.Println(line)
}

func (r dataRow) printRow() {
	line := ""
	for _, measurement := range r.Measurements {
		line = fmt.Sprintf("%s%15s", line, fmt.Sprintf("%f", measurement))
	}
	line = fmt.Sprintf("%s %15v", line, r.IsOutlier)
	fmt.Println(line)
}

func main() {
	rand.Seed(int64(*seed))

	numOutliers := int(float64(*numRows) * (*outlierPercentage))

	clusterCount := make([]int, *dimensionCount)
	clusterOffsets := make(map[int]map[int]float64)
	for i := 0; i < *dimensionCount; i++ {
		clusterCount[i] = *minClusterCount + rand.Intn(*maxClusterCount-*minClusterCount)
		clusterOffsets[i] = make(map[int]float64)
		for j := 0; j < clusterCount[i]; j++ {
			clusterOffsets[i][j] = rand.Float64() * (*clusterMax)
		}
	}

	measurements := make([]dataRow, 0)
	for i := 0; i < *numRows; i++ {
		m := make([]float64, *dimensionCount)
		for j := 0; j < len(m); j++ {
			clusterIdx := rand.Intn(clusterCount[j])
			spread := rand.Float64()*(*clusterSpread) - (*clusterSpread)/2
			m[j] = clusterOffsets[j][clusterIdx] + spread
		}
		measurements = append(measurements, dataRow{
			Measurements: m,
			IsOutlier:    false,
		})
	}

	for i := 0; i < numOutliers; i++ {
		for mIdx := rand.Intn(*numRows); !measurements[mIdx].IsOutlier; {
			measurements[mIdx].IsOutlier = true
			if *multiDimensionalOutliersOnly && !*singleDimensionalOutliersOnly {
				for dim := 0; dim < len(measurements[mIdx].Measurements); dim++ {
					measurements[mIdx].Measurements[dim] = getOutlierValue(clusterOffsets[dim])
				}
			} else if !*multiDimensionalOutliersOnly && *singleDimensionalOutliersOnly {
				dim := rand.Intn(len(measurements[mIdx].Measurements))
				measurements[mIdx].Measurements[dim] = getOutlierValue(clusterOffsets[dim])
			} else {
				addedOutlier := false
				for dim := 0; dim < len(measurements[mIdx].Measurements); dim++ {
					if rand.Intn(2) == 0 {
						continue
					}
					measurements[mIdx].Measurements[dim] = getOutlierValue(clusterOffsets[dim])
					addedOutlier = true
				}
				if !addedOutlier {
					dim := rand.Intn(len(measurements[mIdx].Measurements))
					measurements[mIdx].Measurements[dim] = getOutlierValue(clusterOffsets[dim])
				}
			}
		}
	}

	generatedMeasurements := measurements
	if *shuffle {
		rand.Shuffle(len(generatedMeasurements), func(i, j int) {
			generatedMeasurements[i], generatedMeasurements[j] = generatedMeasurements[j], generatedMeasurements[i]
		})
	}

	if *humanReadableFormat {
		measurements[0].printHeader()
		for _, measurement := range generatedMeasurements {
			measurement.printRow()
		}
	} else {
		measurements[0].printCsvHeader()
		for _, measurement := range generatedMeasurements {
			measurement.printCsvRow()
		}
	}

}

func getOutlierValue(prohibitedRanges map[int]float64) float64 {
	value := prohibitedRanges[0]
	for isProhibited(value, prohibitedRanges) {
		isNegative := rand.Intn(2) == 0 && *generateNegativeOutliers
		value = rand.Float64() * (*clusterMax * 1.25)
		if isNegative {
			value = value * -1
		}
	}
	return value
}

func isProhibited(v float64, prohibitedRanges map[int]float64) bool {
	for _, f := range prohibitedRanges {
		min := f - (*clusterSpread)/2
		max := f + (*clusterSpread)/2
		if v >= min && v <= max {
			return true
		}
	}

	return false
}

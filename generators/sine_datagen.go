package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	numRows           = flag.Int("num-rows", 10000, "The number of rows datasets to generate")
	dimensionCount    = flag.Int("num-dims", 2, "The number of input dimensions (e.g. sensors)")
	outlierPercentage = flag.Float64("outlier-percentage", 0.1, "The percentage (in decimal) of outliers in the output dataset")

	amplitude = flag.Float64("amplitude", 50, "The Amplitude of the sine wave")
	offset    = flag.Float64("offset", 80, "The offset from zero to add to the sine value")

	divider = flag.Float64("divider", 50, "The index divider for the sine wave")

	maxValue                 = flag.Float64("max", 65.0, "Maximum Value")
	outlierDelta             = flag.Float64("outlier-delta", 100.0, "Outlier outage")
	generateNegativeOutliers = flag.Bool("negative-outliers", true, "Generate negative outliers")

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

	startIdxes := make([]int, *dimensionCount)
	for i := 0; i < *dimensionCount; i++ {
		startIdxes[i] = rand.Intn(*numRows)
	}

	measurements := make([]dataRow, 0)
	for i := 0; i < *numRows; i++ {
		m := make([]float64, *dimensionCount)
		for j := 0; j < len(m); j++ {
			sineIdx := (float64(i) / (*divider)) + float64(startIdxes[j])
			m[j] = getSineValue(sineIdx)
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
					measurements[mIdx].Measurements[dim] = getOutlierValue(i)
				}
			} else if !*multiDimensionalOutliersOnly && *singleDimensionalOutliersOnly {
				dim := rand.Intn(len(measurements[mIdx].Measurements))
				measurements[mIdx].Measurements[dim] = getOutlierValue(i)
			} else {
				for dim := 0; dim < len(measurements[mIdx].Measurements); dim++ {
					if rand.Intn(2) == 0 {
						continue
					}
					measurements[mIdx].Measurements[dim] = getOutlierValue(i)
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

func getSineValue(idx float64) float64 {
	return *offset + math.Sin(idx)*(*amplitude)
}

func getOutlierValue(idx int) float64 {
	min := *offset - *amplitude
	max := *offset + *amplitude
	value := *offset + (*amplitude)*0.5
	for value >= min && value <= max {
		isNegative := rand.Intn(2) == 0 && *generateNegativeOutliers
		value = *maxValue + rand.Float64()*(*outlierDelta)
		if isNegative {
			value = value * -1
		}
	}
	return value
}

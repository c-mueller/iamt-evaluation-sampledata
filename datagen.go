package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

/*
This utility can be used to generate data to test the IAMT Validation Engine with

Mode of operation:
The application is invoked by supplying some configuration variables using command line arguments.
Including:


If an argument is not set the default value is used, this one can be determined by

*/

var (
	numRows           = flag.Int("num-rows", 10000, "The number of rows datasets to generate")
	dimensionCount    = flag.Int("num-dims", 2, "The number of input dimensions (e.g. sensors)")
	outlierPercentage = flag.Float64("outlier-percentage", 0.1, "The percentage (in decimal) of outliers in the output dataset")

	minValue                 = flag.Float64("min", 5.0, "Minimum Value")
	maxValue                 = flag.Float64("max", 65.0, "Maximum Value")
	outlierDelta             = flag.Float64("outlier-delta", 100.0, "Outlier outage")
	generateNegativeOutliers = flag.Bool("negative-outliers", true, "Generate negative outliers")

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
	numNonOutliers := *numRows - numOutliers

	measurements := make([]dataRow, 0)
	for i := 0; i < numNonOutliers; i++ {
		m := make([]float64, *dimensionCount)
		for j := 0; j < len(m); j++ {
			m[j] = getValue()
		}
		measurements = append(measurements, dataRow{
			Measurements: m,
			IsOutlier:    false,
		})
	}

	outliers := make([]dataRow, 0)

	for i := 0; i < numOutliers; i++ {
		m := make([]float64, *dimensionCount)
		for j := 0; j < len(m); j++ {
			if *multiDimensionalOutliersOnly && !*singleDimensionalOutliersOnly {
				m[j] = getOutlierValue()
			} else if *singleDimensionalOutliersOnly && !*multiDimensionalOutliersOnly {
				if j == 0 {
					m[j] = getOutlierValue()
				} else {
					m[j] = getValue()
				}
			} else {
				if j == 0 {
					m[j] = getOutlierValue()
				} else {
					if rand.Intn(2) == 0 {
						m[j] = getOutlierValue()
					} else {
						m[j] = getValue()
					}
				}
			}
		}
		if len(m) > 1 {
			rand.Shuffle(len(m), func(i, j int) {
				m[i], m[j] = m[j], m[i]
			})
		}
		outliers = append(outliers, dataRow{
			IsOutlier:    true,
			Measurements: m,
		})
	}

	generatedMeasurements := make([]dataRow, 0)
	generatedMeasurements = append(generatedMeasurements, outliers...)
	generatedMeasurements = append(generatedMeasurements, measurements...)
	rand.Shuffle(len(generatedMeasurements), func(i, j int) {
		generatedMeasurements[i], generatedMeasurements[j] = generatedMeasurements[j], generatedMeasurements[i]
	})

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

func getValue() float64 {
	delta := *maxValue - *minValue
	return *minValue + rand.Float64()*delta
}

func getOutlierValue() float64 {
	isNegative := rand.Intn(2) == 0 && *generateNegativeOutliers
	value := *maxValue + rand.Float64()*(*outlierDelta)
	if isNegative {
		value = value * -1
	}
	return value
}

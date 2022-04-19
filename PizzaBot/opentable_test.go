package main

import (
	"testing"
	"time"
)

func TestQuery(t *testing.T) {

	potentialTime, _ := time.Parse(time.RFC3339, "2019-04-01T18:30:00.000+00:00")
	result, _ := Query("24712", potentialTime, "5")

	print(result.Results)

}

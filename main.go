package main

import "flag"
import log "github.com/sirupsen/logrus"

func main() {
	inputFile := flag.String("inputFile", "./servers.txt", "File with the list of hosts to be processed")
	outputFile := flag.String("outputFile", "./output.txt", "File with the metrics agregated in json format")
	numOfWorkers := flag.Int("numOfWorkers", 4, "Number of concurrent workers that will be issuing http requests")

	flag.Parse()
	err := processInputFile(*inputFile, *outputFile, *numOfWorkers)

	if err != nil {
		log.Error(err.Error())
	}
}

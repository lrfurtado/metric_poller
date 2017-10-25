package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

//Read text file with server list and return them through an output channel
func generateJobs(filePath string) (<-chan string, error) {
	c := make(chan string)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	go func() {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(c)
	}()

	return c, nil
}

func worker(id int, jobs <-chan string, results chan Result, wg *sync.WaitGroup) {
	for j := range jobs {
		log.Info("worker", id, "started  job", j)
		resp, err := http.Get(fmt.Sprintf("http://%v/status", j))
		if err != nil {
			log.Error("worker ", id, " finished job with errors: ", err.Error())
			continue
		}
		res, err := parseResult(resp.Body)
		if err != nil {
			log.Error("worker ", id, " finished job with errors: ", err.Error())
			continue
		}
		results <- *res
		log.Info("worker", id, "finished job", j)
	}

	wg.Done()
}

func processInputFile(inputPath string, outputPath string, numOfWorkers int) error {
	results := make(chan Result)
	jobs, err := generateJobs(inputPath)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		go worker(i, jobs, results, &wg)
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	aggr := make(map[string]Result)

	for r := range results {
		aggrRes := aggr[r.Application]
		aggrRes.Application = r.Application
		aggrRes.RequestCount += r.RequestCount
		aggrRes.SuccessCount += r.SuccessCount
		aggrRes.ErrorCount += r.ErrorCount
		aggrRes.SuccessRate = float64(aggrRes.SuccessCount) / float64(aggrRes.RequestCount)
		aggr[r.Application] = aggrRes
	}

	f, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	aggrList := []Result{}
	data := [][]string{}
	for _, v := range aggr {
		aggrList = append(aggrList, v)
		data = append(data, []string{v.Application, strconv.FormatInt(v.RequestCount, 10), strconv.FormatInt(v.SuccessCount, 10), fmt.Sprint(v.SuccessRate)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Application", "Total Req Count", "Total Succ Count", "Success Rate"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	err = encoder.Encode(aggrList)
	if err != nil {
		return err
	}

	return nil
}

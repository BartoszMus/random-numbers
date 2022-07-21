package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Stddev float64 `json:"stddev"`
	Data   []int   `json:"data"`
}

func GetQuery(w http.ResponseWriter, r *http.Request) (int, int, error) {
	requests, err := strconv.Atoi(r.URL.Query()["requests"][0])
	if err != nil {
		log.Println(err)
		return 0, 0, errors.New("Invalid number of requests")
	}
	length, err := strconv.Atoi(r.URL.Query()["length"][0])
	if err != nil {
		log.Println(err)
		return 0, 0, errors.New("Invalid request length")
	}
	return requests, length, nil
}

func GetNumbersFromWebsite(requests int, length int) ([][]int, error) {
	var chanArray []chan []int
	var numbersTotal [][]int
	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {
		ctx, cancelFunction := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelFunction()
		chanArray = append(chanArray, make(chan []int))

		//Performing requests concurently
		go func(ctx context.Context, ch chan []int) {
			wg.Add(1)
			defer close(ch)
			defer wg.Done()

			//Sending request to the website
			url := fmt.Sprintf("https://www.random.org/integers/?num=%v&min=1&max=10&format=plain&col=1&base=10", length)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				log.Println(err)
				return
			}

			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.Println(resp.StatusCode)
				return
			}
			defer resp.Body.Close()

			//Conversion of the data from website to integer numbers
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
			bodyString := string(bodyBytes)
			stringNumbers := strings.Fields(bodyString)
			numberValues, err := ConvertToNumbers(stringNumbers)
			if err != nil {
				log.Println(err)
				return
			}
			ch <- numberValues
			return
		}(ctx, chanArray[i])
	}

	wg.Wait()

	//Saving all requests in a single array
	for i := 0; i < requests; i++ {
		numbers := <-chanArray[i]
		numbersTotal = append(numbersTotal, numbers)
	}
	return numbersTotal, nil
}

func ConvertToNumbers(stringNumbers []string) ([]int, error) {
	var numberValues []int
	for _, i := range stringNumbers {
		parseInt, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Unable to convert page contents to numbers")
		}
		numberValues = append(numberValues, int(parseInt))
	}
	return numberValues, nil
}

func FormatData(numbers [][]int) ([]Response, error) {
	var allNumbers []int
	var response []Response
	for _, data := range numbers {
		standardDeviation, err := CalculateStandardDeviation(data)
		if err != nil {
			log.Println(err)
			return nil, errors.New("Unable to calculate standard deviation")
		}
		response = append(response, Response{standardDeviation, data})
		allNumbers = append(allNumbers, data...)
	}
	standardDeviationAll, err := CalculateStandardDeviation(allNumbers)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Unable to calculate standard deviation")
	}
	response = append(response, Response{standardDeviationAll, allNumbers})

	return response, nil
}

func CalculateStandardDeviation(values []int) (float64, error) {
	if len(values) == 0 {
		return 0.0, errors.New("No data to calculate standard deviation from")
	}
	var sum, mean float32
	var standardDeviation float64

	sum = 0
	for _, number := range values {
		sum += float32(number)
	}
	mean = float32(sum / float32(len(values)))

	for _, number := range values {
		a := float32(number) - mean
		standardDeviation += math.Pow(float64(a), 2)
	}
	standardDeviation = math.Sqrt(standardDeviation / float64(len(values)))
	return standardDeviation, nil
}

func GetResponses(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Get %v requests of %s numbers", r.URL.Query()["requests"][0], r.URL.Query()["length"][0])

	requests, length, err := GetQuery(rw, r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	numbers, err := GetNumbersFromWebsite(requests, length)
	if err != nil {
		rw.WriteHeader(http.StatusRequestTimeout)
		return
	}

	response, err := FormatData(numbers)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(response)
}

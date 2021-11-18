package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var menu = map[string]string{
	"/min":        "Given list of numbers and a quantifier (how many) provides min number(s)",
	"/max":        "Given list of numbers and a quantifier (how many) provides max number(s)",
	"/avg":        "Given list of numbers calculates their average",
	"/median":     "Given list of numbers calculates their median",
	"/percentile": "Given list of numbers and quantifier 'q', compute the qth percentile of the list elements",
}

type MathRequest struct {
	Numbers []float64
	Quantifier int
}

type MathResponse struct {
	Description string
	Results     []float64
}

func readQueryParam(w http.ResponseWriter, r *http.Request, input string) (string, error) {
	params, ok := r.URL.Query()[input]
	if !ok {
		log.Printf("The query parameter %s is missing\n", input)
		return "", fmt.Errorf("The query param %s is missing", input)
	}
	log.Printf("The query parameter '%s' has value %s\n", input, params)
	return params[0], nil
}

func readNumbersFromUrlQuery(w http.ResponseWriter, r *http.Request) ([]float64, error) {
	numsAsString, err := readQueryParam(w, r, "numbers")
	if err != nil {
		return nil, err
	}
	var nums []float64
	for _, v := range strings.Split(numsAsString, ",") {
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			nums = append(nums, num)
		}
	}
	return nums, nil
}

func readQuantifierFromUrlQuery(w http.ResponseWriter, r *http.Request) (int, error) {
	q, err := readQueryParam(w, r, "quantifier")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(q)
}

func readRequest(w http.ResponseWriter, r *http.Request) MathRequest {
	nums, _ := readNumbersFromUrlQuery(w, r)
	q, _ := readQuantifierFromUrlQuery(w, r)
	return MathRequest{Numbers: nums, Quantifier: q}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }
  if r.Method == "GET" {
    w.Write([]byte("<h1>Welcome to the math-api web server!</h1>"))
	  w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(menu)
  } else {
    http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
  }
}

func Min(w http.ResponseWriter, r *http.Request, request MathRequest) []byte {
	c := Calculation{Numbers: request.Numbers}
	payload, _ := json.Marshal(MathResponse{fmt.Sprintf("Min of %#v with quantifier %v", request.Numbers, request.Quantifier), c.Min(request.Quantifier)})
	return payload
}

func Max(w http.ResponseWriter, r *http.Request, request MathRequest) []byte {
	c := Calculation{Numbers: request.Numbers}
	payload, _ := json.Marshal(MathResponse{fmt.Sprintf("Max of %#v with quantifier %v", request.Numbers, request.Quantifier), c.Max(request.Quantifier)})
	return payload
}

func Avg(w http.ResponseWriter, r *http.Request, request MathRequest) []byte {
	c := Calculation{Numbers: request.Numbers}
	payload, _ := json.Marshal(MathResponse{fmt.Sprintf("Average of %#v", request.Numbers), []float64{c.Average()}})
	return payload
}

func Median(w http.ResponseWriter, r *http.Request, request MathRequest) []byte {
	c := Calculation{Numbers: request.Numbers}
	payload, _ := json.Marshal(MathResponse{fmt.Sprintf("Median of %#v", request.Numbers), []float64{c.Median()}})
	return payload
}

func Percentile(w http.ResponseWriter, r *http.Request, request MathRequest) []byte {
	c := Calculation{Numbers: request.Numbers}
	payload, _ := json.Marshal(MathResponse{fmt.Sprintf("%vth percentile of %#v", request.Quantifier, request.Numbers), []float64{c.Percentile(request.Quantifier)}})
	return payload
}

func Handler(fn func(http.ResponseWriter, *http.Request, MathRequest) []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := readRequest(w, r)
		payload := fn(w, r, request)
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}

func main() {
	mux := http.NewServeMux()
  mux.HandleFunc("/", indexHandler)
	// Using fixed paths
  mux.HandleFunc("/min", Handler(Min))
	mux.HandleFunc("/max", Handler(Max))
	mux.HandleFunc("/avg", Handler(Avg))
	mux.HandleFunc("/median", Handler(Median))
	mux.HandleFunc("/percentile", Handler(Percentile))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMathApiHandlers(t *testing.T) {
	testData := []struct {
		handler      func(w http.ResponseWriter, r *http.Request)
		url          string
		expectedJson string
	}{
		{
			handler:      Handler(Min),
			url:          "/min?numbers=70,30,80,50,90&quantifier=1",
			expectedJson: `{"Description":"Min of []float64{70, 30, 80, 50, 90} with quantifier 1","Results":[30]}`,
		},
		{
			handler:      Handler(Max),
			url:          "/max?numbers=70,30,80,50,90&quantifier=1",
			expectedJson: `{"Description":"Max of []float64{70, 30, 80, 50, 90} with quantifier 1","Results":[90]}`,
		},
		{
			handler:      Handler(Avg),
			url:          "/avg?numbers=70,30,80,50,90",
			expectedJson: `{"Description":"Average of []float64{70, 30, 80, 50, 90}","Results":[64]}`,
		},
		{
			handler:      Handler(Median),
			url:          "/median?numbers=70,30,80,40,50,90",
			expectedJson: `{"Description":"Median of []float64{70, 30, 80, 40, 50, 90}","Results":[60]}`,
		},
		{
			handler:      Handler(Percentile),
			url:          "/percentile?numbers=40,50,60,70,88,44,55,66,4,33,89,90,99,78,76,98,78,95,93,92&quantifier=90",
			expectedJson: `{"Description":"90th percentile of []float64{40, 50, 60, 70, 88, 44, 55, 66, 4, 33, 89, 90, 99, 78, 76, 98, 78, 95, 93, 92}","Results":[95]}`,
		},
	}

	for _, testDatum := range testData {
		t.Run(fmt.Sprintf(" %s returns 200 OK with response: %s", testDatum.url, testDatum.expectedJson), func(t *testing.T) {
			request, err := http.NewRequest("GET", testDatum.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()
			handler := http.HandlerFunc(testDatum.handler)
			handler.ServeHTTP(recorder, request)
			if status := recorder.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
			if recorder.Body.String() != testDatum.expectedJson {
				t.Errorf("handler returned unexpected body: got %v want %v",
					recorder.Body.String(), testDatum.expectedJson)
			}
		})
	}
}

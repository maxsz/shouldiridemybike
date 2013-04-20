package shouldiridemybike

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
)

const API_KEY = ""
const API_URL = "https://api.forecast.io/forecast/" + API_KEY

const DECISION_MIN_TEMP = 5
const DECISION_MAX_TEMP = 30
const DECISION_MAX_PRECIP = 0.20
const DECISION_MAX_WINDSPEED = 15.5

type Page struct {
	Title     string
}

type Decision struct {
	Result bool
	Reason string
	Error *string
}

type Location struct {
	Latitude  string
	Longitude string
}

// forecast.io API objects
type DataPoint struct {
	Time                   float64
	Summary                string
	Icon                   string
	SunriseTime            float64
	PrecipIntensity        float64
	PrecipIntensityMax     float64
	PrecipIntensityMaxTime float64
	PrecipProbability      float64
	PrecipType             string
	PrecipAccumulation     float64
	Temperature            float64
	WindSpeed              float64
	WindBearing            float64
	CloudCover             float64
	Humidity               float64
	Pressure               float64
	Visibility             float64
}

type DataBlock struct {
	Summary string
	Icon    string
	Data []DataPoint
}

type Forecast struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	Offset    float64
	Currently DataPoint
	Minutely  DataBlock
	Hourly    DataBlock
	Daily     DataBlock
}

// Make a decision based on the given forecast
func decide(forecast *Forecast) (*Decision) {
	result := true
	var reason string

	elem := forecast.Hourly.Data[0]
	if elem.Temperature < DECISION_MIN_TEMP || elem.Temperature > DECISION_MAX_TEMP {
		var temp_highlow string
		result = false
		if (elem.Temperature < DECISION_MIN_TEMP) {
			temp_highlow = "low"
		} else {
			temp_highlow = "high"
		}
		temp := strconv.FormatFloat(elem.Temperature, 'f', -1, 64)
		reason += "Temperature too "+ temp_highlow +": " + temp + "°C."
	} else {
		reason += "Temperature is fine."
	}

	if elem.PrecipProbability > DECISION_MAX_PRECIP {
		result = false
		precip := strconv.FormatFloat(elem.PrecipProbability * 100.0, 'f', -1, 64)
		reason += " High chance of rain: " + precip + "%."
	} else {
		reason += " Very low chance of rain."
	}

	if (elem.WindSpeed > DECISION_MAX_WINDSPEED) {
		result = false
		windy := strconv.FormatFloat(elem.WindSpeed * 3.6, 'f', -1, 64)
		reason += " Too windy: " + windy + "km/h."
	} else {
		reason += " Not too windy."
	}

	return &Decision{Result: result, Reason: reason, Error: nil}
}

// Load forecast.io data
func checkForecast(ctx appengine.Context, l Location) (*Forecast, error) {
	var forecast Forecast
	request_uri := API_URL + "/" + l.Latitude + "," + l.Longitude + "?units=si"

	client := urlfetch.Client(ctx)
	resp, err := client.Get(request_uri)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}
	resp.Body.Close()

	err = json.Unmarshal(body, &forecast)

	if nil != err {
		return nil, err
	}

	return &forecast, nil
}

// Handler function for the /shouldi API
func shouldi(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	var decision *Decision
	var encoded_decision []byte

	req.ParseForm()
	form := req.Form
	lat := strings.Join(form["latitude"], "")
	lng := strings.Join(form["longitude"], "")

	f, err := checkForecast(c, Location{lat, lng})
	if err != nil {
		err_string := "Could not acquire forecast data (" + err.Error() + ")"; 
		decision = &Decision{Result: false, Reason: "", Error: &err_string}
		encoded_decision, _ = json.Marshal(decision)
		rw.Write(encoded_decision)
		return
	}
	decision = decide(f)
	encoded_decision, _ = json.Marshal(decision)
	rw.Write(encoded_decision)
}

// Basic index.html template handler
func handler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Should I ride my bike?"}
	t, e := template.ParseFiles("public/index.html")
	if e != nil {
		fmt.Printf("Error: %v\n", e)
		http.Redirect(w, r, "public/notfound.html", http.StatusFound)
		return
	}
	t.Execute(w, p)
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/shouldi", shouldi)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
}

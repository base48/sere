package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/julienschmidt/httprouter"
	bh1750 "github.com/d2r2/go-bh1750"
	i2c "github.com/d2r2/go-i2c"
)

type(
	Data struct{
		Ligh string `json:"ligh"`
		Temp string `json:"temp"`
		Pres string `json:"pres"`
	}
)

func GetSen(){
	i2c, _ := i2c.NewI2C(0x23, 1)
	defer i2c.Close()
	sensor := bh1750.NewBH1750()
	sensor.ChangeSensivityFactor(i2c, 255)
	resolution := bh1750.HighestResolution

	for{
		amb, _ := sensor.MeasureAmbientLight(i2c, resolution)
		datag.Ligh = fmt.Sprintf("%v", amb)
		time.Sleep(1000 * time.Millisecond)
	}
}

var datag Data

func main(){
	r := httprouter.New()
	r.GET("/data", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		data := Data{
			Ligh:	datag.Ligh,
			Temp:	datag.Temp,
			Pres:	datag.Pres,
		}
		dataj, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", dataj)
	})

	go GetSen()								// sensors thread
	http.ListenAndServe("0.0.0.0:3000", r)	// server thread
}

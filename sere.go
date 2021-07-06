package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/julienschmidt/httprouter"
	bh1750 "github.com/d2r2/go-bh1750"
	bsbmp "github.com/d2r2/go-bsbmp"
	i2c "github.com/d2r2/go-i2c"
)

type(
	Data struct{
		Ligh string `json:"illuminance"`
		Temp string `json:"temperature"`
		Pres string `json:"pressure"`
		Humi string `json:"humidity"`
	}
)

func GetSen(){
	i2cl, _ := i2c.NewI2C(0x23, 1) // light sen
	i2co, _ := i2c.NewI2C(0x76, 1) // temp, pres, humi sen
	defer i2cl.Close()
	defer i2co.Close()
	sensorl := bh1750.NewBH1750()
	sensorl.ChangeSensivityFactor(i2cl, 255)
	sensoro, _ := bsbmp.NewBMP(bsbmp.BME280, i2co)

	for{
		a, _ := sensorl.MeasureAmbientLight(i2cl, bh1750.HighestResolution)
		datag.Ligh = fmt.Sprintf("%v", a)
		t, _ := sensoro.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
		datag.Temp = fmt.Sprintf("%v", t)
		p, _ := sensoro.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
		datag.Pres = fmt.Sprintf("%v", p)
		_, h, _ := sensoro.ReadHumidityRH(bsbmp.ACCURACY_STANDARD)
		datag.Humi = fmt.Sprintf("%v", h)
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
			Humi:	datag.Humi,
		}
		dataj, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", dataj)
	})

	go GetSen()								// sensors thread
	http.ListenAndServe("0.0.0.0:10000", r)	// server thread
}

package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"math"

	"github.com/julienschmidt/httprouter"
	bh1750 "github.com/d2r2/go-bh1750"
	bsbmp "github.com/d2r2/go-bsbmp"
	i2c "github.com/d2r2/go-i2c"
)

const(
	ip		= "0.0.0.0"
	port	= 10000
	alt		= 225		// Altitude for pressure correction 
)

type(
	Data struct{
		Ligh uint16 `json:"illuminance"`
		Temp float32 `json:"temperature"`
		Pres float32 `json:"pressure"`
		Humi float32 `json:"humidity"`
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
		datag.Ligh, _ = sensorl.MeasureAmbientLight(i2cl, bh1750.HighestResolution)
		datag.Temp, _ = sensoro.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
		datag.Pres, _ = sensoro.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
		_, datag.Humi, _ = sensoro.ReadHumidityRH(bsbmp.ACCURACY_STANDARD)
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
			Pres:	float32(float64(datag.Pres)/math.Pow(1-alt/44330.0, 5.255)),
			Humi:	datag.Humi,
		}
		dataj, _ := json.Marshal(data)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", dataj)
	})

	go GetSen()	// sensors thread
	http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), r)	// server thread
}

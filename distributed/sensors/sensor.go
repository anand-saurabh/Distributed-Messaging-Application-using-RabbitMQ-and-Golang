package main

// this generates the sensor data to be pushed to rabbitmq.
import
(
  "flag"
  "time"
  "strconv"
  "math/rand"
  "log"
  "encoding/gob"
  "distributed/dto"
  "bytes"
  "distributed/qutils"
  "github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"
var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update the frequency")
var max = flag.Float64("max", 5., "max val for generated readings")
var min = flag.Float64("min", 1., "minimum val of generated readings")
var stepsize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var value = r.Float64()*(*max - * min) + * min
var nom = (*max  - * min)/2 + *min
var r = rand.New(rand.NewSource(time.Now().UnixNano()))


func main()  {

  flag.Parse()

  conn, ch := qutils.GetChannel(url)
  defer conn.Close()
  defer ch.Close()
  dataQueue := qutils.GetQueue(*name, ch)
  dur, _ := time.ParseDuration(strconv.Itoa(1000/(int)(*freq)) + "ms")

  signal := time.Tick(dur)

  buf := new(bytes.Buffer)
  enc := gob.NewEncoder(buf)
  for range signal {
    calcVal()
    reading := dto.SensorMessage{
      Name: *name,
      Value : value,
      Timestamp : time.Now(),
    }

    buf.Reset()
    enc.Encode(reading)

    msg := amqp.Publishing {
      Body : buf.Bytes(),
    }

    ch.Publish(
      "",
      dataQueue.Name,
      false,
      false,
      msg)
    log.Printf("reading sent. value %v \n", value)
  }
}
  func calcVal()  {
    var maxStep, minStep float64
    if value < nom {
      maxStep = *stepsize
      minStep = -1 * *stepsize * (value - *min)/ (nom - *min)
    } else {
      maxStep = *stepsize * (*max - value)/ (*max - nom)
      minStep = -1 * *stepsize
    }
    value += r.Float64( ) * (maxStep - minStep) + minStep
  }

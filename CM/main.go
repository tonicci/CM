package main
import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware/logger"
	"github.com/gofiber/fiber/middleware/pprof"
	"log"
	"os"
	"os/exec"
	"time"
)

func setupRoutes(app *fiber.App) {
	app.Post("/Jmeter/Start", jmeterStart)
	app.Get("/Jmeter/Stop",JmeterStop)

}

func main() {
	var port string
	fmt.Printf("Введите порт, пример ввода данных-\":8072\"")
	fmt.Scan(&port)
	f, _ := os.Create("./logs.txt")
	defer f.Close()

	// Create io.Writer
	w := bufio.NewWriter(f)

	// Flush to file every second
	go func() {
		for {
			w.Flush()
			time.Sleep(1 * time.Second)
		}
	}()

	app := fiber.New()

	file, err := os.OpenFile("./123.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app.Use(logger.New(logger.Config{
		Output: file,
	}))
	app.Use(pprof.New())
	setupRoutes(app)
	app.Listen(port)
}
type Jmeter struct {
	Rph   string `json:"rph"`
	BaselineRampup   string `json:"baseline_rampup"`
	BaselinePercent  string `json:"baseline_percent"`
	BaselineDuration string `json:"baseline_duration"`
	StepRampup       string `json:"step_rampup"`
	StepPercent      string `json:"step_percent"`
	StepDuration     string `json:"step_duration"`
}

func jmeterStart(c *fiber.Ctx) error {
	var requests Jmeter
	data:= c.Body()
	json.Unmarshal(data, &requests)
	text := `threadGroup (имя скрипта),TransactionPerHour (операций в час для 100% профиля), MinPacing(время выполнения операции не более...), LG_Count(для НЕ распределённых тестов=1)
uc_01_Login,0,0,1
uc_02_Payment,`+requests.Rph+`,3,1`
	file, err := os.Create(requests.Rph+".txt")

	if err != nil{
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(text)

	fmt.Println("Done.")


	Scenario:="C:\\Users\\ADikin\\Desktop\\jmeter-ganeles\\JmeterSmartScenario.jmx"
	baseline_rampup:="-Jbaseline_rampup="+requests.BaselineRampup
	baseline_percent:="-Jbaseline_percent="+requests.BaselinePercent
	baseline_duration:="-Jbaseline_duration="+requests.BaselineDuration
	step_rampup:="-Jstep_rampup="+requests.StepRampup
	step_percent:="-Jstep_percent="+requests.StepDuration
	step_duration:="-Jstep_duration="+requests.StepDuration
	profile:="-Jprofile=C:\\Users\\ADikin\\Desktop\\CM_UC\\"+requests.Rph+".txt"
	cmd := exec.Command("C:\\Users\\ADikin\\Desktop\\apache-jmeter-5.4.1\\bin\\jmeter.bat", "-n", "-t",Scenario,profile, baseline_rampup, baseline_percent,baseline_duration,step_rampup,step_percent,step_duration)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)

	//c.Send(response1)
	return nil
}
func JmeterStop(c *fiber.Ctx) error {
	cmd := exec.Command("taskkill", "/f","/im", "java.exe")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
	c.SendString("Stoped ALL test")
	return nil

}



package loggen

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	useragentList []string
	categoryList  []string
	endTime       = time.Now()
	currentTime   = endTime
)

func causeErr(errRate float64) bool {
	rand.Seed(time.Now().UnixNano())
	for i := 1.0; ; i *= 10 {
		if errRate*i >= 1.0 {
			n := rand.Intn(100 * int(i))
			if int(errRate*i) > n {
				return true
			} else {
				return false
			}
		}
	}
}

func floatToIntString(input float64) string {
	return strconv.Itoa(int(input))
}

func randInt(min int, max int) string {
	return strconv.Itoa(min + rand.Intn(max-min))
}

// Generate random number based on log-normal distribution
func randLogNormal(mu, sigma float64) float64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1).Float64()
	r2 := rand.New(s1).Float64()
	z := mu + sigma*math.Sqrt(-2.0*math.Log(r1))*math.Sin(2.0*math.Pi*r2)
	return math.Exp(z)
}

func returnNewList(path string) []string {
	var newList []string

	absPath, _ := filepath.Abs(path)
	fp, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		newList = append(newList, scanner.Text())
	}

	return newList
}

func zitter(i int) int {
	min := i - 3
	max := i + 3
	if min <= 0 {
		return 0
	}
	return min + rand.Intn(max-min)
}

// TODO: exclude private IP address
func Ipv4Address() string {
	var ipStr string
	ipStr = randInt(1, 223) + "."
	ipStr += randInt(0, 255) + "."
	ipStr += randInt(0, 255) + "."
	ipStr += randInt(0, 255)
	return ipStr
}

func RequestTime(i int) string {
	returnTime := currentTime.Add(time.Second * time.Duration(i))
	return returnTime.Format("02/Jan/2006:15:04:05 -0700")
}

func RequestType() string {
	s := []string{"GET", "POST", "PUT", "DELETE"}
	return s[rand.Intn(len(s))]
}

func ReturnStatusCode(errRate float64) string {
	rand.NewSource(time.Now().UnixNano())
	if causeErr(errRate) == false {
		return "200"
	} else {
		s := []string{"301", "403", "404", "500"}
		return s[rand.Intn(len(s))]
	}
}

func ReturnUserAgent() string {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(useragentList) == 0 {
		useragentList = returnNewList(os.Getenv("GOPATH") + "/src/github.com/acroquest/apache-loggen-go/resources/useragents.txt")
	}
	useragent := useragentList[rand.Intn(len(useragentList))]
	return useragent
}

func ReturnRequest() string {
	if len(categoryList) == 0 {
		categoryList = returnNewList(os.Getenv("GOPATH") + "/src/github.com/acroquest/apache-loggen-go/resources/categories.txt")
	}
	category := categoryList[rand.Intn(len(categoryList))]

	i := rand.Intn(10)
	if i < 7 {
		return "\"" + RequestType() + " /category/" + category + " HTTP/1.1\" "
	} else {
		return "\"" + RequestType() + " /" + category + "/" + randInt(1, 999) + " HTTP/1.1\" "
	}
}

func ReturnReferer() string {
	referer := "-"
	return "\"" + referer + "\""
}

func ReturnRecord(i int, errRate float64) string {
	bytes := randInt(20, 5000)
	referer := ReturnReferer()
	responseTime := floatToIntString(20000 * randLogNormal(0.0, 0.5))
	return Ipv4Address() + " - - [" + RequestTime(i) + "] " + ReturnRequest() + ReturnStatusCode(errRate) + " " + bytes + " " + referer + " \"" + ReturnUserAgent() + "\" " + responseTime
}

// TODO change the amount of log data every day.
func GenerateLog(days int, errRate float64) {
	currentTime = endTime.Add(-24 * time.Hour * time.Duration(days))

	// generating log data every 1 second
	for i := 0; endTime.Sub(currentTime.Add(time.Second*time.Duration(i))) >= 0; i += 1 {
		hour := currentTime.Add(time.Second * time.Duration(i)).Hour()
		rand.Seed(time.Now().UnixNano())
		j := rand.Intn(10)

		switch {
		case hour >= 1 && hour <= 5:
			if j <= 2 {
				fmt.Println(ReturnRecord(i, errRate))
			}
		case hour >= 6 && hour <= 9:
			if j <= 4 {
				fmt.Println(ReturnRecord(i, errRate))
			}
		case hour >= 10 && hour <= 17:
			if j <= 6 {
				fmt.Println(ReturnRecord(i, errRate))
			}
		case hour >= 18 && hour <= 23:
			if j <= 8 {
				fmt.Println(ReturnRecord(i, errRate))
			}
		default:
			if j <= 6 {
				fmt.Println(ReturnRecord(i, errRate))
			}
		}
	}
}

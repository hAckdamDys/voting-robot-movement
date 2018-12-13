package main

import (
	"github.com/kataras/iris"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func generateVotes(s *SafeCommands) {
	for {
		s.mux.Lock()
		for command := 0; command < commandSize; command++ {
			valTmp, _ := strconv.Atoi(s.votes[command])
			s.votes[command] = strconv.Itoa(valTmp + rand.Intn(3))
		}
		s.mux.Unlock()
		time.Sleep(time.Second * 1)
	}
}

func resetVotes(s *SafeCommands, waitTime int, stepLoss int) {
	for {
		s.mux.Lock()
		curMax := 0
		newCommand := idle // if no votes then idle
		votes := s.votes
		perm := rand.Perm(commandSize) // we always take random permutation so there are no favourites
		for _, command := range perm {
			valTmp, _ := strconv.Atoi(votes[command])
			if valTmp > curMax {
				curMax = valTmp
				newCommand = Command(command)
			}
		}
		//s.votes = [5]string{"0", "0", "0", "0", "0"}
		idleVotes, _ := strconv.Atoi(votes[idle])
		forVotes, _ := strconv.Atoi(votes[forward])
		backVotes, _ := strconv.Atoi(votes[backward])
		leftVotes, _ := strconv.Atoi(votes[left])
		rightVotes, _ := strconv.Atoi(votes[right])
		s.votes[forward] = strconv.Itoa(int(math.Max(float64(forVotes-stepLoss), 0)))
		s.votes[backward] = strconv.Itoa(int(math.Max(float64(backVotes-stepLoss), 0)))
		s.votes[left] = strconv.Itoa(int(math.Max(float64(leftVotes-stepLoss), 0)))
		s.votes[right] = strconv.Itoa(int(math.Max(float64(rightVotes-stepLoss), 0)))
		s.votes[idle] = strconv.Itoa(int(math.Max(float64(idleVotes-stepLoss), 0)))
		s.lastCommand = String(newCommand)
		s.mux.Unlock()
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}
}

func resetVotesAvg(s *SafeCommands, waitTime int, stepLoss int, multSpeed int, multSteer int, multBrake int) {
	for {
		s.mux.Lock()
		//newCommand := "idle" // if no votes then idle
		leftWheel := 0
		rightWheel := 0
		votes := s.votes

		forMat := [2]int{multSpeed, multSpeed}
		backMat := [2]int{-multSpeed, -multSpeed}
		leftMat := [2]int{multSteer, -multSteer}
		rightMat := [2]int{-multSteer, multSteer}
		idleMat := [2]int{multBrake, multBrake} //decrease by that much for each vote must be positive
		idleVotes, _ := strconv.Atoi(votes[idle])

		forVotes, _ := strconv.Atoi(votes[forward])
		leftWheel += forMat[0] * forVotes
		rightWheel += forMat[1] * forVotes

		backVotes, _ := strconv.Atoi(votes[backward])
		leftWheel += backMat[0] * backVotes
		rightWheel += backMat[1] * backVotes

		leftVotes, _ := strconv.Atoi(votes[left])
		leftWheel += leftMat[0] * leftVotes
		rightWheel += leftMat[1] * leftVotes

		rightVotes, _ := strconv.Atoi(votes[right])
		leftWheel += rightMat[0] * rightVotes
		rightWheel += rightMat[1] * rightVotes

		if leftWheel != 0 {
			if leftWheel > 0 {
				leftWheel = int(math.Max(float64(leftWheel-idleVotes*idleMat[0]), 0))
			} else {
				leftWheel = int(math.Min(float64(leftWheel+idleVotes*idleMat[0]), 0))
			}
		}

		if rightWheel != 0 {
			if rightWheel > 0 {
				rightWheel = int(math.Max(float64(rightWheel-idleVotes*idleMat[1]), 0))
			} else {
				rightWheel = int(math.Min(float64(rightWheel+idleVotes*idleMat[1]), 0))
			}
		}

		//s.votes = [5]string{"0", "0", "0", "0", "0"}
		s.votes[forward] = strconv.Itoa(int(math.Max(float64(forVotes-stepLoss), 0)))
		s.votes[backward] = strconv.Itoa(int(math.Max(float64(backVotes-stepLoss), 0)))
		s.votes[left] = strconv.Itoa(int(math.Max(float64(leftVotes-stepLoss), 0)))
		s.votes[right] = strconv.Itoa(int(math.Max(float64(rightVotes-stepLoss), 0)))
		s.votes[idle] = strconv.Itoa(int(math.Max(float64(idleVotes-stepLoss), 0)))
		s.lastCommand = strconv.Itoa(leftWheel) + "|" + strconv.Itoa(rightWheel)
		s.mux.Unlock()
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
	}
}

type Command int

const publicDir = "public";
const commandSize = 5
const (
	idle     Command = 0
	forward  Command = 1
	backward Command = 2
	left     Command = 3
	right    Command = 4
)

func String(command Command) string {
	names := [...]string{
		"idle",
		"forward",
		"backward",
		"left",
		"right"}
	return names[command]
}
func stringToCommand(s string) Command {
	switch s {
	case "idle":
		return idle
	case "forward":
		return forward
	case "backward":
		return backward
	case "left":
		return left
	case "right":
		return right
	default:
		return -1
	}
}

// safe to use concurrently.
type SafeCommands struct {
	votes       [5]string // string because it is most often used as get
	lastCommand string
	mux         sync.Mutex
}

var page = struct {
	Title string
}{"Vote Fast!!!"}

func main() {
	argsWithoutProg := os.Args[1:]
	port := "8080"
	isAvgMethod := true
	waitTime := 2000
	stepLoss := 10
	multSpeed := 1
	multSteer := 1
	multBrake := 1
	for _, arg := range argsWithoutProg {
		if strings.HasPrefix(arg, "--port") {
			port = strings.Split(arg, "=")[1]
		}else if strings.HasPrefix(arg, "--method") {
			arg = strings.Split(arg, "=")[1]
			if strings.Contains(arg, "avg") {
				isAvgMethod = true
			}
			if strings.Contains(arg, "single") {
				isAvgMethod = false
			}
		}else if strings.HasPrefix(arg, "--wait") {
			arg = strings.Split(arg, "=")[1]
			waitTime, _ = strconv.Atoi(arg)
		}else if strings.HasPrefix(arg, "--steploss") {
			arg = strings.Split(arg, "=")[1]
			stepLoss, _ = strconv.Atoi(arg)
		}else if strings.HasPrefix(arg, "--multspeed") {
			arg = strings.Split(arg, "=")[1]
			multSpeed, _ = strconv.Atoi(arg)
		}else if strings.HasPrefix(arg, "--multsteer") {
			arg = strings.Split(arg, "=")[1]
			multSteer, _ = strconv.Atoi(arg)
		}else if strings.HasPrefix(arg, "--multbrake") {
			arg = strings.Split(arg, "=")[1]
			multBrake, _ = strconv.Atoi(arg)
		}else{
			println("port=what port webserver should run on, default 8080")
			println("method=one method from below")
			println("avg -> default method -> changes votes to (left wheel|right wheel) speed, for example left vote = (-1|1) speed")
			println("forward = (1|1) , backward = (-1|-1)")
			println("left = (-1|1), right = (1|-1)")
			println("idle/brake = (+1 towards 0 value | +1 towards 0 value)")
			println("single -> most voted movement is choosen")
			println("wait=how many miliseconds between reset/decrease of votes, default 2000")
			println("steploss=how many votes decreases per every wait period, default 10")
			println("Multiplication for avg method votes, default 1:")
			println("multspeed=forward,backward votes become (multspeed|multspeed) (-multspeed|-multspeed)")
			println("multsteer=left,right votes become (-multsteer|multsteer) (multsteer|-multsteer)")
			println("multbrake=idle/brake become (+multbrake towards 0 value | +multbrake towards 0 value)")
			println("Example command with all parameters")
			println("go run main.go --port=80 --method=avg --wait=3000 --steploss=30 --multspeed=20 --multsteer=30 --multbrake=50")
			os.Exit(1)
		}
	}

	app := iris.New()
	//app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML(publicDir, ".html"))

	//app.Use(recover.New())
	//app.Use(logger.New())
	var commandsToDo SafeCommands
	if isAvgMethod {
		commandsToDo = SafeCommands{votes: [5]string{"0", "0", "0", "0", "0"}, lastCommand: "0|0"}
	} else {
		commandsToDo = SafeCommands{votes: [5]string{"0", "0", "0", "0", "0"}, lastCommand: "idle"}
	}
	app.StaticWeb("/", publicDir)
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})

	app.Get("/getCommand", func(ctx iris.Context) {
		//println(commandsToDo.lastCommand)
		ctx.HTML(commandsToDo.lastCommand)
	})

	app.Post("/postDirection", func(ctx iris.Context) {
		//s := `form:"direction"`
		var s string
		ctx.ReadForm(&s)
		commandsToDo.mux.Lock()
		command := stringToCommand(s)
		if command != -1 {
			valTmp, _ := strconv.Atoi(commandsToDo.votes[command])
			commandsToDo.votes[command] = strconv.Itoa(valTmp + 1)
		} else {
			println("invalid command: " + s)
		}
		commandsToDo.mux.Unlock()
	})

	app.Handle("GET", "/directions", func(ctx iris.Context) {
		commandsToDo.mux.Lock()
		ctx.HTML(commandsToDo.votes[0] + "|" + commandsToDo.votes[1] + "|" + commandsToDo.votes[2] + "|" + commandsToDo.votes[3] + "|" + commandsToDo.votes[4])
		commandsToDo.mux.Unlock()
	})
	assetHandler := app.StaticHandler(publicDir, false, false)

	app.SPA(assetHandler)

	//go generateVotes(&commandsToDo)
	if isAvgMethod {
		go resetVotesAvg(&commandsToDo, waitTime, stepLoss, multSpeed, multSteer, multBrake)
	} else {
		go resetVotes(&commandsToDo, waitTime, stepLoss)
	}
	app.Run(iris.Addr("0.0.0.0:"+port), iris.WithoutServerError(iris.ErrServerClosed))
}

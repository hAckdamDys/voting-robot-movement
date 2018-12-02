package main

import (
	"github.com/kataras/iris"
	"math/rand"
	"strconv"
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

func resetVotes(s *SafeCommands) {
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
		s.votes = [5]string{"0", "0", "0", "0", "0"}
		s.lastCommand = newCommand
		s.mux.Unlock()
		time.Sleep(time.Second * 10)
	}
}

type Command int

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
	lastCommand Command
	mux         sync.Mutex
}

var page = struct {
	Title string
}{"Vote Fast!!!"}

func main() {
	app := iris.New()
	//app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("public", ".html"))

	//app.Use(recover.New())
	//app.Use(logger.New())

	commandsToDo := SafeCommands{votes: [5]string{"0", "0", "0", "0", "0"}, lastCommand: idle}

	app.StaticWeb("/", "public")
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})

	app.Get("/getCommand", func(ctx iris.Context) {
		//println(commandsToDo.lastCommand)
		ctx.HTML(String(commandsToDo.lastCommand))
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
	assetHandler := app.StaticHandler("public", false, false)

	app.SPA(assetHandler)

	//go generateVotes(&commandsToDo)
	go resetVotes(&commandsToDo)
	app.Run(iris.Addr("0.0.0.0:8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

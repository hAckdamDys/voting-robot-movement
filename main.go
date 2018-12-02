package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

//func say(s chan string) {
//	for ; ;  {
//		time.Sleep(time.Second*5)
//		s <- "forward"
//		time.Sleep(time.Second*5)
//		s <- "backward"
//	}
//}
func generateVotes(s *SafeCommand) {
	for {
		s.mux.Lock()
		valTmp, _ := strconv.Atoi(s.votes[0])
		s.votes[0] = strconv.Itoa(valTmp + rand.Intn(3))
		valTmp, _ = strconv.Atoi(s.votes[1])
		s.votes[1] = strconv.Itoa(valTmp + rand.Intn(3))
		valTmp, _ = strconv.Atoi(s.votes[2])
		s.votes[2] = strconv.Itoa(valTmp + rand.Intn(3))
		valTmp, _ = strconv.Atoi(s.votes[3])
		s.votes[3] = strconv.Itoa(valTmp + rand.Intn(3))
		s.mux.Unlock()
		time.Sleep(time.Second * 1)

	}
}

func resetVotes(s *SafeCommand) {
	for {
		s.mux.Lock()
		s.votes = [4]string{"0", "0", "0", "0"}
		s.mux.Unlock()
		time.Sleep(time.Second * 10)
	}
}

// safe to use concurrently.
type SafeCommand struct {
	votes [4]string
	mux   sync.Mutex
}

var page = struct {
	Title string
}{"Welcome"}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("public", ".html"))
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	//commandToDo := make(chan string)
	commandToDo := SafeCommand{votes: [4]string{"0", "0", "0", "0"}}

	app.StaticWeb("/", "public")
	app.Get("/", func(ctx iris.Context) {
		ctx.ViewData("Page", page)
		ctx.View("index.html")
	})

	app.Handle("GET", "/directions", func(ctx iris.Context) {
		commandToDo.mux.Lock()
		ctx.HTML(commandToDo.votes[0] + "|" + commandToDo.votes[1] + "|" + commandToDo.votes[2] + "|" + commandToDo.votes[3])
		commandToDo.mux.Unlock()
	})
	assetHandler := app.StaticHandler("public", false, false)

	app.SPA(assetHandler)

	go generateVotes(&commandToDo)
	go resetVotes(&commandToDo)
	app.Run(iris.Addr("0.0.0.0:8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

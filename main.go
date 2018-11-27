package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
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
func say(s *SafeCommand) {
	for {
		s.mux.Lock()
		s.direction = "forward"
		s.value = (s.value + 1) % 5
		s.mux.Unlock()
		time.Sleep(time.Second * 3)
		s.mux.Lock()
		s.direction = "backward"
		s.value = (s.value + 1) % 5
		s.mux.Unlock()
		time.Sleep(time.Second * 3)

	}

}

// safe to use concurrently.
type SafeCommand struct {
	direction string
	value     int
	mux       sync.Mutex
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	//commandToDo := make(chan string)
	commandToDo := SafeCommand{direction: "forward", value: 1}
	app.Handle("GET", "/", func(ctx iris.Context) {

		commandToDo.mux.Lock()
		ctx.HTML((commandToDo.direction) + strconv.Itoa(commandToDo.value))
		commandToDo.mux.Unlock()
		//println("xd")
		//time.Sleep(time.Second*5)
		//ctx.HTML("<h1>Welcome 2</h1>")
	})

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	go say(&commandToDo)
	app.Run(iris.Addr("0.0.0.0:8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

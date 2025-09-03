package main

import (
	"blogaggregator/internal/config"
	"blogaggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var logger = log.Default()

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		logger.Print(err)
	}
	newState := State{}
	newState.cfgPtr = &cfg
	db, err := sql.Open("postgres", newState.cfgPtr.DbURL)
	if err != nil {
		fmt.Println(err)
		logger.Print(err)
	}
	dbQueries := database.New(db)
	newState.db = dbQueries
	newCommands := Commands{}
	newCommands.commandList = make(map[string]func(*State, Command) error)
	newCommands.register("login", handlerLogin)
	newCommands.register("register", handlerRegister)
	newCommands.register("reset", handlerReset)
	newCommands.register("users", handlerGetUsers)
	newCommands.register("agg", middlewareLoggedIn(handlerAgg))
	newCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed)) //takes two args
	newCommands.register("feeds", handlerFeeds)
	newCommands.register("follow", middlewareLoggedIn(handlerFollow))
	newCommands.register("following", middlewareLoggedIn(handlerFollowsForUser))
	newCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	newCommands.register("browse", middlewareLoggedIn(handlerBrowse))
	args := os.Args
	if args[1] == "login" || args[1] == "register" {

		if len(args) < 2 {
			fmt.Println("not enough arguments")
			os.Exit(1)
		}
		if len(args) < 3 {
			fmt.Println("username is required")
			os.Exit(1)

		}
	}
	command := Command{}
	command.name = args[1]
	command.args = args[2:]
	err = newCommands.run(&newState, command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		logger.Fatal(err)
		os.Exit(1)
	}

}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, c Command) error {
		user, err := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUsername)
		if err != nil {
			logger.Print(err)
			return err

		}
		err = handler(s, c, user)
		if err != nil {
			logger.Print(err)
			return err
		}
		return nil
	}

}

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"blogaggregator/internal/database"

	"github.com/google/uuid"
)

func handlerLogin(s *State, cmd Command) error {
	ctx := context.Background()
	if len(cmd.args) == 0 {
		return fmt.Errorf("command error: no arguments")
	}
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user doesn't exit")
	}
	err = s.cfgPtr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("message sent")
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	ctx := context.Background()
	userArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err == nil {
		return fmt.Errorf("user already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	newUser, err := s.db.CreateUser(ctx, userArgs)
	if err != nil {
		return err
	}
	err = s.cfgPtr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("user created successfully: \n %v", newUser)
	return nil

}

func handlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	err := s.db.PurgeUsers(ctx)
	if err != nil {
		fmt.Println("error reseting users table")
		return err
	}
	fmt.Println("users table reset successfully")
	return nil
}

func handlerGetUsers(s *State, cmd Command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.cfgPtr.CurrentUsername {
			user = user + " (current)"
		}
		fmt.Println("* " + user)
	}
	return nil
}

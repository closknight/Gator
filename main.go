package main

import (
	"Gator/internal/config"
	"Gator/internal/database"
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading file\n")
	}
	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	programState := &state{cfg: &config, db: dbQueries}
	cmds := commands{table: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		log.Fatalf("Error: arguements needed")
	}

	cm := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(programState, cm)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

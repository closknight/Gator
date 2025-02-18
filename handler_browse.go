package main

import (
	"Gator/internal/database"
	"context"
	"fmt"
	"strconv"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("limit should be an integer")
		}
	} else {
		limit = 2
	}
	posts, err := s.db.GetPostFromUser(context.Background(), database.GetPostFromUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("error retrieving posts: %w", err)
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		fmt.Printf("* %s\n", post.Description.String)
	}
	return nil
}

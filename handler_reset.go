package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting users table: %w", err)
	}
	fmt.Println("users table succcesfully reseted")
	return nil
}

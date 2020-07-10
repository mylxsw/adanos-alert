package repository_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
)

func TestExpectReadyAt(t *testing.T) {
	fmt.Println(repository.ExpectReadyAt("09:00"))
}

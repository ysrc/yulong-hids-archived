package models

import (
	"fmt"
	"testing"
)

func TestDistinct(t *testing.T) {
	cli := NewClient()
	fmt.Println(cli.Distinct("ip"))
}

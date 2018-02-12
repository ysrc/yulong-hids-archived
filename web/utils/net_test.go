package utils

import (
	"fmt"
	"testing"
)

func TestTCPAlive(t *testing.T) {
	fmt.Println(TCPAlive("1222"))
}

func TestString2NetIP(t *testing.T) {
	fmt.Println("[*]", String2NetIP("127.0.0.1"))
}

func TestBetweenIP(t *testing.T) {
	fmt.Println(BetweenIP("127.0.0.1", "127.0.0.200"))
}

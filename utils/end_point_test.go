package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndPoint_String(t *testing.T) {
	ep := EndPoint{IP: "10.10.123.1", Port: 8080}
	s := ep.String()
	fmt.Println(s)
	assert.Equal(t, "10.10.123.1:8080", s)
}

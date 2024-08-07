package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func RandomStr() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprint(rand.Int())
}

func RightPad(text string, length int) string {
	b := strings.Builder{}
	b.WriteString(text)
	for i := len(text); i <= length; i++ {
		b.WriteString(" ")
	}
	return b.String()
}

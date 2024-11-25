package main

import (
	"fmt"
	"os"

	"github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/transport"
)

func main() {
	fmt.Println(os.Getwd())
	transport.HandleCreate()
}

package cmd

import (
	"fmt"
	"time"
)

func MockRunFun() int {
	fmt.Println("Hello Go!")

	fmt.Println("Mocking some longer operation... cya in 2 seconds!")
	time.Sleep(time.Second * 2)
	return 0
}

func MockCleanupFunc() {
	fmt.Println("Mock cleanup.")
}

uses the log package to log like normal but also save the printed logs to a .txt file in a time dated folder

example usage:

```
package main

import (
	"github.com/Scrimzay/loglogger"
	"path/filepath"
	"os"
)

func main() {
	log, err := logger.New("details.log")
	if err != nil {
		panic(err)
	}
	defer log.Close()

	// basic logger usage
	log.Print("***Program starting***")
	log.Print("This is a log message")
	log.Printf("This is a formatted message: %d", 42)

	str1 := log.Sprint("Hello from Sprint")
	str2 := log.Sprintf("Hello %s", "from Sprintf")
	log.Print(str1)
	log.Print(str2)

	log.Error("This is an error")
	log.Errorf("This is a formatted %s", "error")

	// these will write to the log file and then exit the program
	//log.Fatal("Fatal error occurred")
	log.Fatalf("Fatal error with value: %v", err)
}
```

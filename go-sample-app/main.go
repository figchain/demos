package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-sample-app/figchain"

	"github.com/figchain/go-client/pkg/client"
	"github.com/figchain/go-client/pkg/evaluation"
)

func main() {
	once := flag.Bool("once", false, "Exit after initial fetch")
	flag.Parse()

	log.Println("Starting FigChain Go Sample App...")

	// 1. Initialize Client using config file generated from UI
	// The "client-config.json" file should be placed in the current directory or provide full path
	// Double check the name of the file as the UI will generate the file with the namespace in the name.
	configPath := "client-config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s not found. Please place the generated config file in the root.", configPath)
	}

	// We can also override specific settings if needed, using a second argument to NewClientFromConfig
	c, err := client.NewClientFromConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to create FigChain client: %v", err)
	}
	defer c.Close()

	log.Println("Client initialized successfully.")

	// 2. Define the Fig Key we are interested in.
	// Be sure to match the key name with a key that you defined in the UI
	// Also, be sure that your schema matches the field that you access in this code.
	// For example, if you defined a key named "test" in the UI, and schema with a string field named "test", then no code changes are required in this sample app.

	figKey := "test" // Replace with actual key if different
	log.Printf("Listening for updates on key: %s", figKey)

	// 3. Register Callback
	// We use a prototype instance of the generated struct.
	// The FigChain Go generator places models in the 'avro' package
	// Assuming the generated struct is named 'Test'

	c.RegisterListener(figKey, &figchain.Test{}, func(r client.AvroRecord) {
		testRecord, ok := r.(*figchain.Test)
		if !ok {
			log.Printf("Received record is not of type *figchain.Test")
			return
		}

		// Pretty print the updated config
		log.Printf(">>> UPDATE RECEIVED for %s <<<", figKey)
		log.Printf("New Configuration: test=%+v", testRecord.Test)
	})

	// Fetch initial value
	evalContext := evaluation.NewEvaluationContext(nil)
	var initialVal figchain.Test
	if err := c.GetFig(figKey, &initialVal, evalContext); err != nil {
		log.Printf("Initial GetFig failed: %v", err)
	} else {
		log.Printf("Initial value fetched. test=%+v", initialVal.Test)
	}

	if *once {
		log.Println("Exiting because -once was specified.")
		return
	}

	// Keep the application running
	log.Println("Waiting for updates... (Press Ctrl+C to exit)")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Shutting down...")
}

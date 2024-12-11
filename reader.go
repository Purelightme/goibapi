package ibapi

import (
	"bufio"
	"context"
	"sync"
)

// EReader starts the scan and decode goroutines
func EReader(ctx context.Context, client *EClient, scanner *bufio.Scanner, decoder *EDecoder, wg *sync.WaitGroup) {

	msgChan := make(chan []byte, 300)

	// Decode
	wg.Add(1)
	go func() {
		log.Debug().Msg("decoder started")
		defer log.Debug().Msg("decoder ended")
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				decoder.interpret(msg) // single worker and no go here!!
			}
		}
	}()

	// Scan
	wg.Add(1)
	go func() {
		log.Debug().Msg("scanner started")
		defer client.Disconnect()
		defer log.Debug().Msg("scanner ended")
		defer wg.Done()
		for {
			hasMore := scanner.Scan()
			if hasMore {
				msgBytes := make([]byte, len(scanner.Bytes()))
				copy(msgBytes, scanner.Bytes())
				msgChan <- msgBytes
			} else {
				err := scanner.Err()
				log.Error().Err(err).Msg("scanner error")
				break
			}
		}
		close(msgChan)
	}()
}

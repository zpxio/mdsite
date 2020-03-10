/*
 * Copyright 2020 zpxio (Jeff Sharpe)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/zpxio/mdsite/pkg/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetHandler(text.New(os.Stderr))

	log.Infof("Starting up...")

	conf := config.Create()
	conf.Load()

	// Set up signal monitoring
	termSignals := make(chan os.Signal, 1)
	exitChan := make(chan bool, 1)
	signal.Notify(termSignals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		shutdownSignal := <-termSignals
		// Clear the output buffer.
		//    ^ Super excessive... but it clears the line buffer if a ^C is printed and lets the logs be pretty.
		fmt.Println()

		log.Infof("Received shutdown signal: %s", shutdownSignal)
		exitChan <- true
	}()

	// Wait for shutdown signals
	<-exitChan
	log.Info("Initiating shutdown.")
}

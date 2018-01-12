// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"bufio"
	"strconv"
	"strings"

	"github.com/XiaoYang.Code4Fun/bitcoin-stats/client"
	dt "github.com/XiaoYang.Code4Fun/bitcoin-stats/data"
)

func readStdin(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSuffix(input, "\n")
}

func waitForStdin(reader *bufio.Reader, done chan bool, c *client.Client) string {
	inputChan := make(chan string)
	go func() {
		inputChan <- readStdin(reader)
	}()
	for {
		select {
		case input := <-inputChan:
			return input
		case <-done:
			cleanupAndExit(c)
		}
	}
}

func generateCSV(start, end int64, dest string, saveCSV chan bool,  c *client.Client) {
	data := make(chan string)
	f, err := os.Create(dest)
	if err != nil {
		return
	}
	w := bufio.NewWriter(f)
	w.WriteString(dt.CSVHeader())
	go func() {
		b, err := c.GetBlockByHeight(start)
		if err != nil {
			fmt.Errorf("Unable to get block at height %v", start)
			saveCSV <- true
			return
		}
		data <- b.CSVBlockOutput()
		for {
			h := b.GetHeight()
			if h >=end {
				break
			}
			b, err = c.GetBlockByHash(b.GetNextBlockHash())
			if err != nil {
				fmt.Errorf("Unable to get block at height %v", h)
				saveCSV <- true
				return
			}
			data <- b.CSVBlockOutput()
		}
		saveCSV <- true
	}()
	for {
		select {
		case d := <-data:
			w.WriteString(d)
		case <-saveCSV:
			w.Flush()
			return
		}
	}
}

func cleanupAndExit(c *client.Client) {
	c.Close()
	os.Exit(1)
}

func main() {
	c, err := client.NewClient("user", "pass")
	if err != nil {
		fmt.Errorf("Unable to establish connection to btcd server %v\n", err)
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	saveCSV := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Terminating the program...")
		done <- true
		saveCSV <- true
	}()

	for {
		if len(done) != 0 {
			cleanupAndExit(c)
		}
		fmt.Println()
		fmt.Println("*********************************************")
		fmt.Println("* Please select from the following options: *")
		fmt.Println("* 1 Show current block count.               *")
		fmt.Println("* 2 Read a block.                           *")
		fmt.Println("* 3 Generate csv for a range of blocks.     *")
		fmt.Println("*********************************************")
		reader := bufio.NewReader(os.Stdin)
		input := waitForStdin(reader, done, c)
		switch (input) {
		case "1":
			blockCount, err := c.GetBlockCount()
			if err != nil {
				fmt.Errorf("Unable to obtain current block count %v\n", err)
			}
			fmt.Printf("Current block count is %v\n", blockCount)
		case "2":
			fmt.Print("Please enter the block height or hash: ")
			blockId := waitForStdin(reader, done, c)
			blockHeight, err := strconv.Atoi(blockId)
			if err != nil {
				b, err := c.GetBlockByHash(blockId)
				if err != nil {
					fmt.Errorf("Cannot get block by its hash %v", err)
					cleanupAndExit(c)
				}
				fmt.Printf("Block at hash:\n%v\n", b.SingleBlockOutput())
			} else {
				b, err := c.GetBlockByHeight(int64(blockHeight))
				if err != nil {
					fmt.Errorf("Cannot get block by its height %v", err)
					cleanupAndExit(c)
				}
				fmt.Printf("Block at height:\n%v\n", b.SingleBlockOutput())
			}
		case "3":
			fmt.Println("Please enter the starting block height, ending block height, and destination csv file (seperated by space): ")
			sed := waitForStdin(reader, done, c)
			s := strings.Split(sed, " ")
			if len(s) != 3 {
				fmt.Printf("Missing parameters in %v", sed)
			}
			start := strings.TrimSpace(s[0])
			end := strings.TrimSpace(s[1])
			dest := strings.TrimSpace(s[2])
			sHeight, err := strconv.Atoi(start)
			if err != nil {
				fmt.Errorf("Block height must be an integter, %v", err)
				continue
			}
			eHeight, err := strconv.Atoi(end)
			if err != nil {
				fmt.Errorf("Block height must be an integter, %v", err)
				continue
			}
			fmt.Printf("Generating csv for blocks from %v to %v to %v\n", sHeight, eHeight, dest)
			generateCSV(int64(sHeight), int64(eHeight), dest, saveCSV, c)
		default:
			fmt.Errorf("Unknown option: %v\n", input)
		}
	}
}

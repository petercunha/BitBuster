# BitBuster
A cross-platform Instagram account cracker written in Go


## How to Use
1. Run `go get -u -v github.com/petercunha/goinsta` to download dependencies.

2. You need two files in the same directory as BitBuster.
  * `passwords.txt`: Needs to contain the passwords you want to try against the Instagram account.
  * `proxies.txt`: Needs to contain working and checked HTTPS proxies. You can use a tool [like this](https://github.com/chill117/proxy-lists) to scrape them.

3. Once you have these files ready to go, simply edit the "USERNAME" constant (near the top of `main.go`), and change it to the person's Instagram username who you'd like to hack.

4. Finally, open up a terminal and run `go run main.go`

Now you can sit back, relax, and let BitBuster do the hard work!

package main

import (
	"fmt"
	"github.com/joffarex/ytptk-go/googleauth"
)

func main() {
	googleauth.AuthenticateWithGoogle()
	fmt.Println("FUCK YEAH!")
}

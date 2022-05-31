package utils

import (
	"crypto/rand"
	"log"
	"os"
)

//gen salt of length n
func GenSalt(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetPfp(u string) string {
	pfp_def := os.Getenv("PFPS")
	pfp_def = "/pfps/"

	if pfp_def == "" {
		log.Fatal("can't get pfp path")
	}

	pfp := pfp_def + u + ".jpeg"
	_, err := os.Stat(pfp)

	if err == nil {
		return pfp
	}

	return pfp_def + "guest.jpeg"
}

package service

import (
	"fmt"
	"math/rand"
	"time"
)

func generateOTP() string {
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%04d", randGenerator.Intn(10000))
}

const otpCooldown = 50 * time.Second

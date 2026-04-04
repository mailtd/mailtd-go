package mailtd

import (
	"crypto/sha256"
	"strconv"
	"time"
)

const defaultDifficulty = 15

// PoWSolution represents a proof-of-work solution for account creation.
type PoWSolution struct {
	T     int64  `json:"t"`
	N     string `json:"n"`
	D     int    `json:"d"`
	Token string `json:"token,omitempty"`
}

// SolvePow computes a proof-of-work for the given address and difficulty.
func SolvePow(address string, difficulty int) PoWSolution {
	timestamp := time.Now().Unix()
	tsStr := strconv.FormatInt(timestamp, 10)

	for nonce := 0; ; nonce++ {
		nonceStr := strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(address + tsStr + nonceStr))
		if hasLeadingZeroBits(hash[:], difficulty) {
			return PoWSolution{
				T: timestamp,
				N: nonceStr,
				D: difficulty,
			}
		}
	}
}

func hasLeadingZeroBits(hash []byte, bits int) bool {
	remaining := bits
	for _, b := range hash {
		if remaining <= 0 {
			return true
		}
		if remaining >= 8 {
			if b != 0 {
				return false
			}
			remaining -= 8
		} else {
			if b>>(8-remaining) != 0 {
				return false
			}
			return true
		}
	}
	return remaining <= 0
}

// powRetryResponse is the server response when a PoW step-up is required.
type powRetryResponse struct {
	Status             string `json:"status"`
	RequiredDifficulty int    `json:"required_difficulty"`
	Token              string `json:"token"`
}

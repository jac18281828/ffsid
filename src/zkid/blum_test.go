package main

import (
	"math/big"
	"testing"
)

func TestSecureRandom(t *testing.T) {
	bitrange := [8]uint{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	for i, _ := range bitrange {
		r, err := SecureRandom(bitrange[i])
		if err != nil {
			t.Error(err)
		}
		if r.BitLen() != int(bitrange[i]) {
			t.Errorf("Expected %d, got %d", bitrange[i], r.BitLen())
		}
	}
}

func TestIsBlumInteger(t *testing.T) {
	blum := [8]string{
		"3",
		"7",
		"11",
		"15",
		"19",
		"23",
		"27",
		"31",
	}
	for i, _ := range blum {
		bigBlum := new(big.Int)
		bigBlum.SetString(blum[i], 10)
		if !IsBlumPrime(bigBlum) {
			t.Errorf("Expected %s to be a Blum Integer", blum[i])
		}
	}
}

func TestChooseBlumInteger(t *testing.T) {
	for i := 0; i < 10; i++ {
		b, p, q, err := ChooseBlumInteger(256)
		if err != nil {
			t.Error(err)
		}
		q1 := new(big.Int).Div(b, p)
		if q1.Cmp(q) != 0 {
			t.Errorf("Expected %d, got %d", q, q1)
		}
		if !IsBlumPrime(p) {
			t.Errorf("Expected %d to be a Blum Prime", p)
		}
		if !IsBlumPrime(q) {
			t.Errorf("Expected %d to be a Blum Prime", q)
		}
	}
}

func TestEasyToFactor(t *testing.T) {
	N := int64(256)
	for i := int64(2); i < N; i++ {
		factor := big.NewInt(i)
		isPrime := IsPrime(factor, N)
		easyToFactor := EasyToFactor(factor)
		println(i, isPrime, easyToFactor)
		if isPrime && easyToFactor {
			if isPrime {
				t.Errorf("Expected %d to be prime", i)
			} else {
				t.Errorf("Expected %d to factor", i)
			}
		}
	}
}

func IsPrime(n *big.Int, k int64) bool {
	j := new(big.Int)
	if j.Mod(n, big.NewInt(2)).Int64() == 0 {
		return false
	}
	for i := int64(3); i < k; i += 1 {

		if j.Mod(n, big.NewInt(i)).Int64() == 0 {
			return false
		}
	}
	return true
}

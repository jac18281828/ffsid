package main

import (
	"crypto/rand"
	"math/big"
)

const MILLER_RABIN_ROUNDS = 20

func IsBlumPrime(n *big.Int) bool {
	n = new(big.Int).Set(n)
	THREE := big.NewInt(3)
	FOUR := big.NewInt(4)
	isCongruent3Mod4 := n.Mod(n, FOUR).Cmp(THREE) == 0
	isPrime := n.ProbablyPrime(MILLER_RABIN_ROUNDS)
	return isCongruent3Mod4 && isPrime
}

func ChooseBlumPrime(bits uint) (*big.Int, error) {
	r, err := SecureRandom(bits)
	if err != nil {
		return big.NewInt(0), err
	}
	for !IsBlumPrime(r) {
		r, err = SecureRandom(bits)
		if err != nil {
			return big.NewInt(0), err
		}
	}
	return r, nil
}

func ChooseBlumInteger(bits uint) (b *big.Int, p *big.Int, q *big.Int, err error) {
	b = big.NewInt(4)
	p = big.NewInt(2)
	q = big.NewInt(2)
	for EasyToFactor(b) || p.Cmp(q) == 0 || new(big.Int).Sub(p, q).BitLen() < int(bits/4) {
		p, err = ChooseBlumPrime(bits / 2)
		if err != nil {
			NIL := big.NewInt(0)
			return NIL, NIL, NIL, err
		}
		q, err = ChooseBlumPrime(bits / 2)
		if err != nil {
			NIL := big.NewInt(0)
			return NIL, NIL, NIL, err
		}
		b.Mul(p, q)
	}
	return b, p, q, nil

}

func SecureRandom(bits uint) (*big.Int, error) {
	r := big.NewInt(1)
	err := error(nil)
	max := big.NewInt(1).Lsh(r, bits+1)
	for r.BitLen() < int(bits) || r.BitLen() > int(bits) {
		r, err = rand.Int(rand.Reader, max)
		if err != nil {
			return big.NewInt(0), err
		}
		r.SetBit(r, int(bits-1), 1)
		r.SetBit(r, 0, 1)
	}
	return r, nil
}

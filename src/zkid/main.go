package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// avoids modulus that is trivial to factor
func EasyToFactor(n *big.Int) bool {
	j := new(big.Int).Set(n)
	if j.Mod(n, big.NewInt(2)).Int64() == 0 {
		return true
	}
	for i := int64(3); i < 256; i += 2 {
		if j.Mod(n, big.NewInt(i)).Int64() == 0 {
			return true
		}
	}
	return false
}

func ChooseQuadraticResidue(p *big.Int, bits uint) (*big.Int, error) {
	x, err := SecureRandom(bits)
	if err != nil {
		return big.NewInt(0), err
	}
	r := new(big.Int).Exp(x, big.NewInt(2), p)
	for r.Cmp(p) >= 0 {
		x, err = SecureRandom(bits)
		if err != nil {
			return big.NewInt(0), err
		}
		r = new(big.Int).Exp(x, big.NewInt(2), p)
	}
	return r, nil
}

func GeneratePrivateKey(k int, bits uint, n *big.Int) []*big.Int {
	/**
	 * generate private key which i s a list of k random quadratic residues of n
	 */
	var privateKey []*big.Int
	for i := 0; i < k; {
		r, err := ChooseQuadraticResidue(n, bits)
		if err != nil {
			continue
		}
		for _, s := range privateKey {
			if s.Cmp(r) == 0 {
				continue
			}
		}
		privateKey = append(privateKey, r)
		i++
	}
	return privateKey
}

func GeneratePublicKey(privateKey []*big.Int, n *big.Int) []*big.Int {
	var publicKey []*big.Int
	for _, r := range privateKey {
		modInverseTerm := new(big.Int).Exp(r, big.NewInt(-2), n)
		if modInverseTerm != nil {
			publicKey = append(publicKey, modInverseTerm)
		} else {
			panic("modInverseTerm is nil")
		}
	}
	return publicKey
}

func ArrayToString(a []*big.Int) string {
	s := ""
	for _, v := range a {
		s += v.String()
		s += " "
	}
	return s
}

func PeggyWitness(privateKey []*big.Int, R *big.Int, challenge *big.Int, n *big.Int) *big.Int {
	witness := big.NewInt(1)
	witness.Mul(witness, R)
	challenge_value := challenge.Uint64()
	for _, s := range privateKey {
		last_bit := challenge_value & 1
		witness.Mul(witness, new(big.Int).Exp(s, big.NewInt(int64(last_bit)), n))
		challenge_value = challenge_value >> 1
	}
	return witness
}

func VictorProof(publicKey []*big.Int, Y *big.Int, challenge *big.Int, n *big.Int) *big.Int {
	proof := big.NewInt(1)
	proof.Exp(Y, big.NewInt(2), n)
	challenge_value := challenge.Uint64()
	for _, y := range publicKey {
		last_bit := challenge_value & 1
		proof.Mul(proof, new(big.Int).Exp(y, big.NewInt(int64(last_bit)), n)).Mod(proof, n)
		challenge_value = challenge_value >> 1
	}
	return proof
}

func main() {

	// n will not be a secret and factorization only needs to be valid on
	// the timescale of the proof, values >= 1024 bits are recommended
	const nbits = 1024
	const k_round = 16
	const t_round = 4

	n, _, _, err := ChooseBlumInteger(nbits)

	if err != nil {
		panic(err)
	}
	privateKey := GeneratePrivateKey(k_round, nbits, n)
	publicKey := GeneratePublicKey(privateKey, n)
	println("public modulus n: ", n.String())
	println("public key: ", ArrayToString(publicKey))

	proof_valid := big.NewInt(1)
	round_count := uint(0)

	for i := 0; i < t_round; i++ {
		// accrediditation rounds
		R_peggy, err := SecureRandom(nbits)
		if err != nil {
			panic(err)
		}

		// Victor sees Peggy's R squared
		X := new(big.Int).Exp(R_peggy, big.NewInt(2), n)
		challenge_vector, err := rand.Int(rand.Reader, big.NewInt(1<<k_round))
		Y := PeggyWitness(privateKey, R_peggy, challenge_vector, n)

		victor_expect := VictorProof(publicKey, Y, challenge_vector, n)
		if victor_expect.Cmp(X) != 0 {
			fmt.Println("X: ", X.String())
			fmt.Println("expect: ", victor_expect.String())
			panic("Peggy's proof is invalid")
		}
		round_count += k_round
		proof_valid = proof_valid.Lsh(proof_valid, k_round)
	}
	// probability = 2^(k*t)
	if proof_valid.Cmp(big.NewInt(1<<20)) < 0 {
		panic("Victor can not verify the proof")
	}
	fmt.Println("Peggy and Victor agree with confidence 2 to", round_count, "or 1 in", proof_valid.String())
}

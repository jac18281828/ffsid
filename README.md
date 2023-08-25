# Zero Knowledge Identification 

This is a ZK proof playground for testing identity proofs using the Fiat Shamir method

One party can establish its identity to another by connecting to a trusted mediator. The mediator provides a composite number `n = pq`. Factoring this large number into its prime constituents, p and q, is challenging. The proof is carried out when the prover, Peggy, can demonstrate her knowledge of the quadratic residuosity of n, a task that is difficult without prior knowledge of n's factorization, to the verifier, Victor, in such a way that he can be convinced with an extremely high degree of mathematical certainty.

Protocol Key Generation

1. Peggy chooses k random numbers S1,...Sk modulo the `n` agreed upon from the mediator
2. Peggy reveals her public key Ij consisting of `+-1/Sj^2 (mod n)` - the squared inverse mod n

Efficient proof - accreditation rounds

1. Peggy chooses a random number R and sends `X = +- R^2 (mod n)`
2. Victor sends a challenge vector in the form of a binary string `b1...bk`
3. Peggy replys with `Y = R (s1^b1 * s2^b2 * ... * sk^bk) (mod n)`
4. Victor confirms that `X = +- Y^2 (I1^b1 * I2^b2 * ... * Ik^bk) (mod n)`
5. accreditation proceeds up to t times

Victor is convinced with probability 2 ^ (k * t).  k*t >= 20 is suggested

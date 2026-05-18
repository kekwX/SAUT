# 🛡️ SAUT (haShing, Authentication and aUthenticated encrypTion)

*Read this in other languages: [English](README.md), [Русский](README_ru.md).*

<img src="https://img.shields.io/badge/Go-1.18%2B-00ADD8?style=flat&logo=go" alt="Go Version">

![Architecture](https://img.shields.io/badge/Architecture-Sponge-blueviolet)
![State Size](https://img.shields.io/badge/State-512_bit-success)
![Status](https://img.shields.io/badge/Status-Experimental-orange)

This project provides a Go implementation of the custom cryptographic algorithm **SAUT**. 

Based on the **Sponge construction** (architecturally similar to the Ascon algorithm), SAUT provides Authenticated Encryption with Associated Data (AEAD) and generates a 256-bit hash value (MAC tag) for data integrity verification.

## 🧮 Mathematical Foundation

The internal state of the algorithm is represented as a **4×4** square matrix, where each element is a 32-bit word (yielding a total state size of 512 bits).

Base operations are performed in the Galois Field $GF(2^{32})/P(x)$ using the following irreducible polynomial:
`P(x) = 0x101002005`

The round function (F-function) consists of two main transformations:
1. **Linear Layer:** Matrix multiplication of the state by a constant MDS matrix.
2. **Non-linear Layer:** Application of a custom 4-bit substitution box (S-box) to the matrix columns (implemented via bitwise operations / bit-slicing).

**SAUT Algorithm S-box:**
| x | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 |
|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|:---|
| **S(x)**| 2 | 14 | 13 | 6 | 8 | 10 | 11 | 1 | 5 | 3 | 4 | 9 | 0 | 15 | 7 | 12 |

## ⚙️ Algorithm Phases

The encryption and hashing process is divided into 4 key phases (according to the core specification):

### 1. Initialization
The initial state matrix is filled with zeros.
* The 128-bit constant (IV) is split into 32-bit words and XORed into the **third column** of the matrix. The $F^5$ transformation (5 rounds) is then applied.
* The 128-bit Key is XORed into the **fourth row** of the matrix, followed by another $F^5$ transformation (5 rounds).

### 2. Associated Data Processing (AEAD)
* Associated data is divided into 128-bit blocks. Incomplete blocks are padded using the `10*` rule (a single `1` bit followed by `0`s).
* Each block is XORed into the **second row** of the matrix. After each block, the $F^4$ transformation (4 rounds) is executed.
* Once all AEAD blocks are processed, an additional blank $F^7$ transformation (7 rounds) is applied.

### 3. Encryption (Plaintext Processing)
* The plaintext is divided into 128-bit blocks (with the same `10*` padding rule).
* Each block is XORed with the **second row** of the matrix. The resulting values of this row are extracted as the **ciphertext** ($C_i$).
* After each block, the $F^6$ transformation (6 rounds) is applied.

### 4. Finalization (Hash Generation)
To produce the final 256-bit tag (hash), a squeezing phase is performed:
* The first row of the matrix is temporarily saved ($T$). The $F^{16}$ transformation (16 rounds) is executed.
* The first half of the hash ($H_1$) is computed by XORing the third row with $T$.
* The new first row is saved again ($T'$). The $F^{17}$ transformation (17 rounds) is executed.
* The second half of the hash ($H_2$) is computed by XORing the new third row with $T'$.
* The final hash is the concatenation: $H_1 || H_2$ (256 bits).

## 🚀 Quick Start

### Prerequisites
* [Go](https://golang.org/dl/) installed (version 1.18 or higher).

### Running the Project

1. Clone the repository:
   ```bash
   git clone https://github.com/kekwX/SAUT.git
   cd <FOLDER_NAME>
   ```

2. Run the encryption algorithm (test vectors from the documentation are hardcoded in `main.go`):
   ```bash
   go run main.go
   ```

3. Upon successful execution, a `ciphertext_output.bin` file containing the encrypted data will be generated in the root directory, and the execution time will be printed in the console.

## 📁 Source Code Structure

Key functions in `main.go` mapping to the algorithm's pseudocode:

| Function | Description |
| :--- | :--- |
| `gf32Mul` | Multiplication of two elements in $GF(2^{32})$ with reduction modulo `0x101002005`. |
| `gfADD` | Addition (XOR) in the Galois Field. |
| `mathMulgf` | Multiplication of the state matrix by the constant MDS matrix `M_matrix`. |
| `S_block` | Application of the non-linear S-box to the entire matrix using bitwise logic. |
| `padded` | Implementation of the padding logic to ensure data fits exactly into 128-bit blocks (`10*` padding). |

## 🧪 Test Vectors (Reference)
The code utilizes the following parameters to verify correctness against the specification:
* **IV:** `da58ec5eb389b253dcdffc9e91e00665`
* **Key:** `000102030405060708090a0b0c0d0e0f`
* **AEAD:** `0102030405060708090a0b0c0d0e0f`
* **Plaintext:** `62c256ab2bfe966bd7279c05b22938da`

## ⚠️ Disclaimer
> **Warning:** This implementation is intended strictly for educational, academic, and research purposes to verify the logic of the SAUT algorithm. The code has not undergone a formal security audit and **must not** be used to protect sensitive data in real-world production environments.

---
*Built with ❤️ in Go.*

***


package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const P uint64 = 0x101002005

// два элемента
func gf32Mul(a, b uint32) uint32 {
	var res uint64 = 0
	var test uint32 = 0

	for i := 0; i < 32; i++ {
		test = (b >> i)
		test = test & 1
		if test != 0 {
			res ^= uint64(a) << i
		}
	}
	// Редукция по модулю
	for i := 63; i >= 32; i-- {
		if (res>>i)&1 != 0 {
			res ^= P << (i - 32)
		}
	}
	return uint32(res & 0xFFFFFFFF)
}

// xor двух элементов
func gfADD(a, b uint32) uint32 {
	return a ^ b
}

// две матрицы
func mathMulgf(matrix, M_matrix [4][4]uint32) [4][4]uint32 {
	var C [4][4]uint32
	n := len(matrix)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var sum uint32 = 0
			for k := 0; k < n; k++ {
				prod := uint32(gf32Mul(matrix[i][k], M_matrix[k][j]))
				sum = gfADD(sum, prod)
			}
			C[i][j] = sum
		}
	}
	return C
}

func S_block(new_Lmatrix [4][4]uint32) [4][4]uint32 {

	var buf uint32 = 0x0
	var input uint32 = 0x0
	var output uint32 = 0x0

	for j := 0; j < 4; j++ {
		for k := 31; k >= 0; k-- {
			input = 0
			for i := 0; i < 4; i++ {
				buf = (new_Lmatrix[i][j] >> k) & 1
				input ^= buf << (3 - i)
			}
			//фунцкия S
			output = S_block_change(input)
			//возврат в матрицу
			for v := 0; v < 4; v++ {
				buf = (output >> (3 - v)) & 1
				new_Lmatrix[v][j] &= ^(uint32(1) << k)
				new_Lmatrix[v][j] |= buf << k
			}
		}

	}
	return new_Lmatrix
}

func S_block_change(input uint32) uint32 {
	var S map[uint8]uint8
	var value uint8 = 0

	S = map[uint8]uint8{
		0:  2,
		1:  14,
		2:  13,
		3:  6,
		4:  8,
		5:  10,
		6:  11,
		7:  1,
		8:  5,
		9:  3,
		10: 4,
		11: 9,
		12: 0,
		13: 15,
		14: 7,
		15: 12,
	}

	value = uint8(input & 0xF)
	return uint32(S[value])
}

func padded(stroka string) []uint32 {
	var length int = 0
	var hexStr string = ""
	var arr []uint32

	length = len(stroka)
	length_1 := length % 8
	length_2 := length / 8

	for i := 0; i < length_2; i++ {
		hexStr = ""
		for j := 0; j < 8; j++ {
			hexStr += string(stroka[i*8+j])
		}
		val, _ := strconv.ParseUint(hexStr, 16, 32)
		arr = append(arr, uint32(val))
	}

	if length_1 > 0 {
		hexStr := ""
		for i := length_2 * 8; i < len(stroka); i++ {
			hexStr += string(stroka[i])
		}

		val, _ := strconv.ParseUint(hexStr, 16, 32)
		usedBits := len(hexStr) * 4

		// Добавить 1 бит '1'
		val = (val << 1) | 1
		usedBits++

		// Добавить остальные нули
		if usedBits < 32 {
			val = val << (32 - usedBits)
		}

		arr = append(arr, uint32(val))
	}

	return arr
}

func main() {
	start := time.Now() // фиксируем время начала
	var count int = 1
	var finalCiphertext []uint32
	var outputBytes []byte

	M_matrix := [4][4]uint32{
		{0x02010103, 0x01010302, 0x01030201, 0x03020101},
		{0x03020101, 0x02010103, 0x01010302, 0x01030201},
		{0x01030201, 0x03020101, 0x02010103, 0x01010302},
		{0x01010302, 0x01030201, 0x03020101, 0x02010103},
	}

	// Инициализация матрицы, заполненной нулями
	matrix := [4][4]uint32{
		{0x0, 0x0, 0x0, 0x0},
		{0x0, 0x0, 0x0, 0x0},
		{0x0, 0x0, 0x0, 0x0},
		{0x0, 0x0, 0x0, 0x0},
	}

	// Начальынй вектор
	vector := [4]uint32{0xda58ec5e, 0xb389b253, 0xdcdffc9e, 0x91e00665}

	//ключ
	key := [4]uint32{0x00010203, 0x04050607, 0x08090a0b, 0x0c0d0e0f}

	//AEAD
	// AEAD := [4]uint32{0x01020304, 0x05060708, 0x090a0b0c, 0x0d0e0f}
	AEAD := "0102030405060708090a0b0c0d0e0f"
	AEAD = strings.Repeat(AEAD, 1)
	//plaintext
	//plaintext := [4]uint32{0x62c256ab, 0x2bfe966b, 0xd7279c05, 0xb22938da}
	plaintext := "62c256ab2bfe966bd7279c05b22938da"
	plaintext = strings.Repeat(plaintext, 1024)
	// Запись вектора в третий столбец
	for i := 0; i < 1; i++ {

		for i := 0; i < 4; i++ {
			matrix[i][2] = vector[i]
		}

		new_Lmatrix := mathMulgf(matrix, M_matrix)

		// S-блок
		new_matrix := S_block(new_Lmatrix)

		//остальные 4 раунда после добавления вектора
		for i := 0; i < 4; i++ {
			new_matrix = mathMulgf(new_matrix, M_matrix)
			new_matrix = S_block(new_matrix)
		}

		//добавление ключа
		for i := 0; i < 4; i++ {
			new_matrix[3][i] ^= key[i]
		}

		//5 раундов после добавления ключа
		for i := 0; i < 5; i++ {
			new_matrix = mathMulgf(new_matrix, M_matrix)
			new_matrix = S_block(new_matrix)
		}

		//AEAD
		count = len(AEAD) / 32
		if count == 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			NEW_AEAD := padded(AEAD)
			for i := 0; i < 4; i++ {
				new_matrix[1][i] ^= NEW_AEAD[i]
			}
			for i := 0; i < 4; i++ {
				new_matrix = mathMulgf(new_matrix, M_matrix)
				new_matrix = S_block(new_matrix)
			}
		}

		//7 раундов после добавления (последнего)AEAD
		for i := 0; i < 7; i++ {
			new_matrix = mathMulgf(new_matrix, M_matrix)
			new_matrix = S_block(new_matrix)
		}

		//добавление plaintext
		count = len(plaintext) / 32
		if count == 0 {
			count = 1
		}
		// Итоговый шифротекст — с размером count * 4 (по 4 uint32 на блок)
		// finalCiphertext := make([]uint32, count*4)
		for i := 0; i < count; i++ {
			NEW_plaintext := padded(plaintext)
			for i := 0; i < 4; i++ {
				new_matrix[1][i] ^= NEW_plaintext[i]
			}
			finalCiphertext = append(finalCiphertext, new_matrix[1][0], new_matrix[1][1], new_matrix[1][2], new_matrix[1][3])

			for i := 0; i < 6; i++ {
				new_matrix = mathMulgf(new_matrix, M_matrix)
				new_matrix = S_block(new_matrix)
			}

		}
		//сохраняем первую сторку до 16 раундов
		first_row := new_matrix[0]

		//16 раундов после добавления (последнего)plaintext
		for i := 0; i < 22; i++ {
			new_matrix = mathMulgf(new_matrix, M_matrix)
			new_matrix = S_block(new_matrix)
		}

		//сохраняем третью строку после 16 раундов
		third_row := new_matrix[2]

		//xor двух строк и вывод - это первая часть хэша
		hash := make([]uint32, 4)
		for i := 0; i < 4; i++ {
			hash[i] = first_row[i] ^ third_row[i]
		}

		//сохраняем первую строку до 17 раундов
		first_row = new_matrix[0]

		//17 раундов после вывода первой части хэша
		for i := 0; i < 17; i++ {
			new_matrix = mathMulgf(new_matrix, M_matrix)
			new_matrix = S_block(new_matrix)
		}

		//сохраняем третью строку после 17 раундов
		third_row = new_matrix[2]
		//xor двух строк и вывод - это вторая часть хэша
		hash_1 := make([]uint32, 4)
		for i := 0; i < 4; i++ {
			hash_1[i] = first_row[i] ^ third_row[i]
		}

		//коннкатенация двух хэшей
		hash_sum := make([]uint32, 8)
		for i := 0; i < 4; i++ {
			hash_sum[i] = hash[i]
			hash_sum[i+4] = hash_1[i]
		}

		// вывод шифротекста
		// fmt.Printf("Ciphertext ")
		// fmt.Println()
		// for _, val := range finalCiphertext {
		// 	fmt.Printf("%08x ", val)
		// }
		// fmt.Println()

		// вывод хэша
		// fmt.Printf("Hash-sum ")
		// fmt.Println()
		// for _, val := range hash_sum {
		// 	fmt.Printf("%08x ", val)
		// }
		// fmt.Println()

	}
	for _, val := range finalCiphertext {
		// Каждое uint32 (4 байта) преобразуем в байты

		outputBytes = append(outputBytes, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))

	}

	// Записываем байты в файл
	outputFilename := "ciphertext_output.bin"
	err := os.WriteFile(outputFilename, outputBytes, 0644) // 0644 - права доступа
	if err != nil {
		fmt.Printf("Ошибка записи в файл %s: %v\n", outputFilename, err)
	} else {
		fmt.Printf("Шифротекст успешно записан в %s\n", outputFilename)
	}
	duration := time.Since(start) // вычисляем длительность
	seconds := duration.Seconds()
	fmt.Printf("Время выполнения: %.7f секунд\n", seconds)
}

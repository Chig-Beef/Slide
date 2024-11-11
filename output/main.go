package main

func main() {
	for i := 2; i < 100; i++ {
		prime := true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				prime = false
			}
		}
		if prime {
			print(i)
			print()
		}
	}

}

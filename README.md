# gommand
Go one liner program, similar to python -c

Added to the code p() as alias for fmt.Println()

Usage:

    gommand 'fmt.Println("Hello, Gommand!")'
    gommand 'h := md5.New(); io.WriteString(h, "Md5 me"); fmt.Printf("%x", h.Sum(nil))'
    gommand 'p("Array:", []int{1, 2, 3})'

# gommand
Go one liner program, similar to python -c

Added to the code p() as alias for fmt.Println()
If "github.com/k0kubun/pp" is installed - add pp() as alias for pp.Print()

Usage:

    gommand 'fmt.Println("Hello, Gommand!")'
    gommand 'h := md5.New(); io.WriteString(h, "Md5 me"); fmt.Printf("%x", h.Sum(nil))'
    gommand 'p("Array:", []int{1, 2, 3})'
    gommand 'pp(os.Environ())' # get os env, need install "github.com/k0kubun/pp"

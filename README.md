# gommand
Go one liner program, similar to python -c

<strong>Usage</strong>
<ul>
       <li>gommand 'fmt.Println("Hello, Gommand!")' <br /></li>
       <li>gommand 'h := md5.New(); io.WriteString(h, "Md5 me"); fmt.Printf("%x", h.Sum(nil))'</li>
       <li>gommand 'p("Array:", []int{1, 2, 3})'</li>
       <li>gommand 'pp(os.Environ())' # dump os env, need install "github.com/k0kubun/pp"</li>
       <li>gommand 'p(os.Args[1:])' 1 2 3</li>
</ul>

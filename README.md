gommand
=======

Go one liner program, similar to python -c

How to get it?
-------------
```bash
go get github.com/sno6/gommand
```

How to run it?
-------------
```bash
gommand [code]
```

Usage
-----
```bash
gommand 'fmt.Println("Hello, World!")'
```
You can quickly write and run code without worrying about setting up a go file.
gommand auto imports whatever packages are being used by the program so you don't have to worry about it.

Write data to a new file.
```bash
gommand 'f, _ := os.Create("file"); f.Write([]byte("hi")); f.Close()'
```

Run a quick http server on port 8080.
```bash
gommand 'http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "hi") }); http.ListenAndServe(":8080",nil)'
```


Run a quick http server to serve the current directory. Print any errors encountered, such as trying to serve on a port that's already in use.
```bash
gommand 'http.Handle("/", http.FileServer(http.Dir("."))); fmt.Println(http.ListenAndServe(":8080",nil))'
```

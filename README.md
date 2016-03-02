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
You can quickly write and run code without worrying about setting up a file etc.
```bash
       gommand 'f, _ := os.Create("file"); f.Write([]byte("hi")); f.Close()'
```

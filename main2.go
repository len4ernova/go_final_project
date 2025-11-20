package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)
func main2() {
	//ex1()
	fmt.Println("ex2")
	ex2()
}

func ex1() {
	d := http.Dir(".")
	f, err := d.Open("web/css/style.css")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
}

func ex2() {
	d := http.Dir("./web")
	f, err := d.Open("css/style.css")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
}
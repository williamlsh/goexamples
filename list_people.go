package main

import (
	"fmt"
	"goexamples/pb"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:  %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fname := os.Args[1]

	in, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	book := &pb.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	listPeople(os.Stdout, book)
}

func writePerson(w io.Writer, p *pb.Person) {
	fmt.Fprintln(w, "Person ID:", p.Id)
	fmt.Fprintln(w, " Name:", p.Name)
	if p.Email != "" {
		fmt.Fprintln(w, " Email address:", p.Email)
	}
	for _, pn := range p.Phones {
		switch pn.Type {
		case pb.Person_MOBILE:
			fmt.Fprint(w, "  Mobile phone #: ")
		case pb.Person_HOME:
			fmt.Fprint(w, "  Home phone #: ")
		case pb.Person_WORK:
			fmt.Fprint(w, "  Work phone #: ")
		}
		fmt.Println(w, pn.Number)
	}
}

func listPeople(w io.Writer, book *pb.AddressBook) {
	for _, p := range book.People {
		writePerson(w, p)
	}
}

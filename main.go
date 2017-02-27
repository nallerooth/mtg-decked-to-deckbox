package main

// Created by: Nalle Rooth <nalle.rooth [at] gmail [dot] com>

/*
Source format

00: Total Qty,
01: Reg Qty,
02: Foil Qty,
03: Card,
04: Set,
05: Mana Cost,
06: Card Type,
07: Color,
08: Rarity,
09: Mvid,
10: Single Price,
11: Single Foil Price,
12: Total Price,
13: Price Source,
14: Notes

Destination format

00: Count,
01: Tradelist Count,
02: Name,
03: Edition,
04: Card Number,
05: Condition,
06: Language,
07: Foil,
08: Signed,
09: Artist Proof,
10 :Altered Art,
11: Misprint,
12: Promo,
13: Textless,
14: My Price

*/

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type card struct {
	regQty  int
	foilQty int
	number  int
	name    string
	set     string
}

func loadCollectionFromSourceFile(fname string) []*card {
	var err error

	src, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	r := csv.NewReader(src)
	srcCollection, err := r.ReadAll()

	// Allocate storage for collection, minus headers
	collection := make([]*card, len(srcCollection)-1)

	for i, sc := range srcCollection {
		if i == 0 {
			continue
		}

		var c card

		if c.regQty, err = strconv.Atoi(sc[1]); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if c.foilQty, err = strconv.Atoi(sc[2]); err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		if c.number, err = strconv.Atoi(sc[9]); err != nil {
			fmt.Println("Card:", sc)
			fmt.Println("Error: ", err)
		}

		if strings.HasSuffix(sc[4], " Edition") {
			c.set = strings.Replace(sc[4], " Edition", "", 1)
		} else {
			c.set = sc[4]
		}

		c.name = sc[3]

		// Minus 1 because of stripped headers
		collection[i-1] = &c
	}

	return collection
}

func writeDeckboxCSVFile(collection []*card, blacklist map[int]bool) bool {
	dstFile := "converted.csv"

	f, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	headers := []string{
		"Count",
		"Tradelist Count",
		"Name",
		"Edition",
		"Card Number",
		"Condition",
		"Language",
		"Foil",
		"Signed",
		"Artist Proof",
		"Altered Art",
		"Misprint",
		"Promo",
		"Textless",
		"My Price"}

	w.Write(headers)

	tmpl := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		""}

	for _, c := range collection {
		if blacklist[c.number] {
			continue
		}

		tmpl[0] = strconv.Itoa(c.regQty)
		tmpl[2] = c.name
		tmpl[3] = c.set
		tmpl[4] = strconv.Itoa(c.number)
		tmpl[7] = ""

		w.Write(tmpl)

		if c.foilQty > 0 {
			tmpl[0] = strconv.Itoa(c.foilQty)
			tmpl[7] = "1"

			w.Write(tmpl)
		}
	}

	w.Flush()

	return true

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: convert <srcfile>")
		os.Exit(0)
	}

	// Decked Builder adds the back of some 2-sided cards as well, ignore those
	blacklist := map[int]bool{
		409869: true,
		414429: true}

	cards := loadCollectionFromSourceFile(os.Args[1])
	writeDeckboxCSVFile(cards, blacklist)

	fmt.Printf("\nConverted %d unique cards\n", len(cards))
}

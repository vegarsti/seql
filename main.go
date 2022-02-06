package main

import (
	"fmt"
	"seql/spc"
)

func main() {
	r := spc.NewRelation(
		[]string{"name", "from", "resides"},
		[]spc.Row{
			{"Jordan", "New York", "New York"},
			{"Lauren", "California", "New York"},
			{"Justin", "Ontario", "New York"},
			{"Devin", "California", "California"},
			{"Smudge", "Ontario", "Ontario"},
		},
	)
	fmt.Println(r)
	fmt.Println("Lives in New York:")
	fmt.Println(spc.ConstantSelect(r, 2, "New York"))
	fmt.Println("Lives where they're from:")
	fmt.Println(spc.EqualsSelect(r, 1, 2))
	fmt.Println("Only name and resides:")
	fmt.Println(spc.Project(r, []int{0, 2}))
}

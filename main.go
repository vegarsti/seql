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

	fmt.Println("A (big) cross product:")
	c := spc.NewRelation(
		[]string{"location", "country"},
		[]spc.Row{
			{"New York", "United States"},
			{"California", "United States"},
			{"Ontario", "Canada"},
		},
	)
	fmt.Println(spc.Cross(r, c))

	fmt.Println("What country Smudge lives in:")
	fmt.Println(
		// Only keep the country column (1).
		spc.Project(
			// Only keep the row with name (0) = "Smudge".
			spc.ConstantSelect(
				// Throw away everything except the name (0) and the country (4).
				spc.Project(
					// We only want the rows where the "resides" (2) location matches the
					// "location" (3).
					spc.EqualsSelect(
						// First, grab every pair of rows.
						spc.Cross(r, c),
						2,
						3,
					),
					[]int{0, 4},
				),
				0,
				"Smudge",
			),
			[]int{1},
		),
	)
}

package main

import (
	"fmt"
	"seql/spc"
	"seql/volcano"
)

func spcMain() {
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

func volcanoMain() {
	r := volcano.NewRelation(
		[]string{"name", "from", "resides"},
		[]volcano.Row{
			{"Jordan", "New York", "New York"},
			{"Lauren", "California", "New York"},
			{"Justin", "Ontario", "New York"},
			{"Devin", "California", "California"},
			{"Smudge", "Ontario", "Ontario"},
		},
	)
	iter := volcano.ScanRelation(r)
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	fmt.Println("Justin:")
	iter = volcano.ConstantSelect(
		volcano.ScanRelation(r),
		0,
		"Justin",
	)
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	fmt.Println("Lives where they're from:")
	iter = volcano.EqualsSelect(
		volcano.ScanRelation(r),
		1,
		2,
	)
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	fmt.Println("Only name and resides:")
	iter = volcano.Project(volcano.ScanRelation(r), []int{0, 2})
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	c := volcano.NewRelation(
		[]string{"location", "country"},
		[]volcano.Row{
			{"New York", "United States"},
			{"California", "United States"},
			{"Ontario", "Canada"},
		},
	)

	fmt.Println()
	fmt.Println("What country Smudge lives in:")
	iter = volcano.Project(
		// Only keep the row with name (0) = "Smudge".
		volcano.ConstantSelect(
			// Throw away everything except the name (0) and the country (4).
			volcano.Project(
				// We only want the rows where the "resides" (2) location matches the
				// "location" (3).
				volcano.EqualsSelect(
					// First, grab every pair of rows.
					volcano.Cross(
						volcano.ScanRelation(r),
						volcano.ScanRelation(c),
					),
					2,
					3,
				),
				[]int{0, 4},
			),
			0,
			"Smudge",
		),
		[]int{1},
	)
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	// Exercise 1: Union
	fmt.Println("Union")
	iter = volcano.Union(volcano.ScanRelation(c), volcano.ScanRelation(c))
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	// Exercise 2: Zip
	fmt.Println("Zip")
	iter = volcano.Zip(volcano.ScanRelation(r), volcano.ScanRelation(c))
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	// Exercise 3: Inspect
	fmt.Println("Inspect")
	iter = volcano.Inspect(volcano.ScanRelation(c))
	iter.Start()
	for _, ok := iter.Next(); ok; _, ok = iter.Next() {
		// inspect prints itself
	}

	// Exercise 4: Intersect
	fmt.Println("Intersect: Languages in use at companies c1 and c2")
	c1 := volcano.NewRelation(
		[]string{"language"},
		[]volcano.Row{
			{"Go"},
			{"Kotlin"},
		},
	)
	c2 := volcano.NewRelation(
		[]string{"language"},
		[]volcano.Row{
			{"Python"},
			{"Go"},
			{"JavaScript"},
		},
	)
	iter = volcano.Intersect(volcano.ScanRelation(c1), volcano.ScanRelation(c2))
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	// Exercise 5: Distinct
	fmt.Println("Distinct languages used at companies c1 and c2")
	iter = volcano.Distinct(
		volcano.Union(volcano.ScanRelation(c1), volcano.ScanRelation(c2)),
	)
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}

	// Exercise 6 (mine): Order
	fmt.Println()
	fmt.Println("Sort languages used at company c2 using Order")
	iter = volcano.Order(volcano.ScanRelation(c2), []int{0})
	iter.Start()
	for row, ok := iter.Next(); ok; row, ok = iter.Next() {
		fmt.Println(row)
	}
}

func main() {
	volcanoMain()
}

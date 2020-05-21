package tracedeq_test

import (
	"fmt"

	"github.com/nasa9084/tracedeq"
)

type T struct{}

func (t *T) Errorf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

var t T

func Example() {
	type Employee struct {
		Age int
	}
	type Team struct {
		Members map[string]Employee
	}
	type Department struct {
		Teams []Team
	}
	type Company struct {
		Departments []Department
	}

	got := Company{
		Departments: []Department{
			{
				Teams: []Team{
					{
						Members: map[string]Employee{
							"Alice": {
								Age: 25,
							},
							"Bob": {
								Age: 30,
							},
						},
					},
				},
			},
		},
	}
	want := Company{
		Departments: []Department{
			{
				Teams: []Team{
					{
						Members: map[string]Employee{
							"Alice": {
								Age: 20,
							},
							"Bob": {
								Age: 30,
							},
						},
					},
				},
			},
		},
	}

	if result := tracedeq.DeepEqual(got, want); !result.IsEqual {
		t.Errorf("unexpected: %s\n  got:  %v\n  want: %v", result.Trace.Join("."), result.X, result.Y)
		return
	}

	// Output:
	// unexpected: Departments.0.Teams.0.Members.Alice.Age
	//   got:  25
	//   want: 20
}

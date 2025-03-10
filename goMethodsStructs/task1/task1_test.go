package task1

import (
	"testing"
)

func TestPersonGettersWithTableDriven(t *testing.T) {
	testCases := []struct {
		name            string
		person          Person
		expectedName    string
		expectedAge     int
		expectedAddress string
	}{
		{
			name:            "Regular person",
			person:          NewPerson(33, "Rajesh", "Pune"),
			expectedName:    "Rajesh",
			expectedAge:     33,
			expectedAddress: "Pune",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.person.GetName() != tc.expectedName {
				t.Errorf("GetName() returned %s, expected %s", tc.person.GetName(), tc.expectedName)
			}
			if tc.person.GetAge() != tc.expectedAge {
				t.Errorf("GetAge() return %d, expected %d", tc.person.GetAge(), tc.expectedAge)
			}
			if tc.person.GetAddress() != tc.expectedAddress {
				t.Errorf("GetAddress() returned %s, expected %s", tc.person.GetAddress(), tc.expectedAddress)
			}
		})
	}
}

func TestPersonSettersWithTableDriven(t *testing.T) {
	testCases := []struct {
		name            string
		initalPerson    Person
		newAge          int
		newName         string
		newAddress      string
		expectedAge     int
		expectedAddress string
		expectedName    string
	}{
		{
			name:            "Regular test",
			initalPerson:    NewPerson(33, "Raj", "Mumbai"),
			newAge:          34,
			newName:         "Rajeshk",
			newAddress:      "Pune city",
			expectedAge:     34,
			expectedName:    "Rajeshk",
			expectedAddress: "Pune city",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			person := tc.initalPerson

			person.SetName(tc.newName)
			person.SetAge(tc.newAge)
			person.SetAddress(tc.newAddress)

			if person.name != tc.expectedName {
				t.Errorf("After Setname %s, GetName() is %s, expected is %s", tc.newName, person.name, tc.expectedName)
			}
			if person.age != tc.expectedAge {
				t.Errorf("After Setname %d, GetAge() is %d, expected is %d", tc.newAge, person.age, tc.expectedAge)
			}
			if person.address != tc.expectedAddress {
				t.Errorf("After Setaddress %s, GetAddress() is %s, expected is %s", tc.newAddress, person.address, tc.expectedAddress)
			}
		})
	}
}

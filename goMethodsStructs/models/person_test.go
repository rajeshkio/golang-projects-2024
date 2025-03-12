package models

import (
	"testing"
)

func TestPersonGettersWithTableDriven(t *testing.T) {

	address := NewAddress("Maharashtra", "Pune", 411001)
	testCases := []struct {
		name            string
		person          PersonWithAddress
		expectedName    string
		expectedAge     int
		expectedAddress Address
	}{
		{
			name:            "Regular person",
			person:          NewPersonWithAddress(33, "Rajesh", address),
			expectedName:    "Rajesh",
			expectedAge:     33,
			expectedAddress: address,
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
			if tc.person.GetCity() != tc.expectedAddress.GetCity() {
				t.Errorf("GetCity() returned %s, expected %s", tc.person.GetCity(), tc.expectedAddress.GetCity())
			}
		})
	}
}

func TestPersonSettersWithTableDriven(t *testing.T) {
	address := NewAddress("Maharashtra", "Mumbai", 411091)
	testCases := []struct {
		name            string
		initalPerson    PersonWithAddress
		newAge          int
		newName         string
		newAddress      Address
		expectedAge     int
		expectedAddress Address
		expectedName    string
	}{
		{
			name:            "Regular test",
			initalPerson:    NewPersonWithAddress(33, "Raj", address),
			newAge:          34,
			newName:         "Rajeshk",
			newAddress:      address,
			expectedAge:     34,
			expectedName:    "Rajeshk",
			expectedAddress: address,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			person := tc.initalPerson

			person.SetName(tc.newName)
			person.SetAge(tc.newAge)
			person.SetAddress(tc.newAddress)

			if person.name != tc.expectedName {
				t.Errorf("After Setname %s, GetName() is %s, expected is %s", tc.newName, person.GetName(), tc.expectedName)
			}
			if person.age != tc.expectedAge {
				t.Errorf("After Setname %d, GetAge() is %d, expected is %d", tc.newAge, person.GetAge(), tc.expectedAge)
			}
			if person.Address != tc.expectedAddress {
				t.Errorf("After Setaddress %s, GetAddress() is %s, expected is %s", tc.newAddress, person.GetAddress(), tc.expectedAddress)
			}
		})
	}
}

func TestPersonStringWithTableDriven(t *testing.T) {

	address := NewAddress("Maharashtra", "Mumbai", 411091)
	testCases := []struct {
		name           string
		person         PersonWithAddress
		expectedString string
	}{
		{
			name:           "Regular person",
			person:         NewPersonWithAddress(35, "rajesh", address),
			expectedString: "Person(name: rajesh, age: 35, address: Maharashtra, Mumbai, 411091)",
		},
		{
			name:           "Empty fields",
			person:         NewPersonWithAddress(0, "", Address{}),
			expectedString: "Person(name: , age: 0, address: , , 0)",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.person.String() != tc.expectedString {
				t.Errorf("String returned %q, expected %q", tc.person.String(), tc.expectedString)
			}
		})
	}
}

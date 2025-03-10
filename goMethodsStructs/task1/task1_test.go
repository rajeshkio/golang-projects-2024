package task1

import (
	"testing"
)

func TestNewPerson(t *testing.T) {
	person := NewPerson(35, "Anjana", "Bangalore")

	if person.name != "Anjana" {
		t.Errorf("GetName() returned %s, expected 'Anjana'", person.GetName())
	}

	if person.age != 35 {
		t.Errorf("GetAge() returned %d, expected 35", person.GetAge())
	}

	if person.address != "Bangalore" {
		t.Errorf("GetAddress() returned %s, expected 'Bangalore'", person.GetAddress())
	}
}

func TestGettersAndSetters(t *testing.T) {
	person := NewPerson(31, "Krithika", "Pune")

	if person.GetName() != "Krithika" {
		t.Errorf("GetName() returned %s, expected 'Krithika'", person.GetName())
	}

	person.SetName("Rajesh")
	if person.GetName() != "Rajesh" {
		t.Errorf("GetName() returned %s, expected 'Rajesh'", person.GetName())
	}
}

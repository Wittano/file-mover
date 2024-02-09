package filesystem

import "testing"

func TestMountedList(t *testing.T) {
	res, err := MountedList()
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("Mounted list is empty")
	}

	elem := res[0]

	if elem.MountedPoint == "" || elem.Device == "" || elem.Type == "" {
		t.Fatalf("Sample element of mounted list is invalid. MountedPoint: %s, Device: %s, Type: %s", elem.MountedPoint, elem.Device, elem.Type)
	}
}

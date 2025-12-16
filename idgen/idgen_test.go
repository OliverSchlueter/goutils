package idgen

import (
	"testing"
)

func TestGenerateID_Length(t *testing.T) {
	lengths := []int{1, 5, 10, 32}

	for _, l := range lengths {
		id := GenerateID(l)

		if len(id) != l {
			t.Fatalf("expected length %d, got %d", l, len(id))
		}
	}
}

func TestGenerateID_Alphabet(t *testing.T) {
	id := GenerateID(100)

	for i := 0; i < len(id); i++ {
		if !contains(alphabet, id[i]) {
			t.Fatalf("invalid character %q in id %q", id[i], id)
		}
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	for i := 0; i < 5; i++ {
		testUniqueness(t, 8, 1_000_000)
	}
}

//func TestGenerateID_Uniqueness_Enhanced(t *testing.T) {
//	for i := 0; i < 1_000; i++ {
//		testUniqueness(t, 8, 1_000_000)
//	}
//}

func testUniqueness(t *testing.T, idLen, idCount int) {
	seen := make(map[string]struct{}, idCount)

	for i := 0; i < idCount; i++ {
		id := GenerateID(idLen)

		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate id %q generated for length %d", id, idLen)
		}
		seen[id] = struct{}{}
	}
}

//func TestGenerateID_Uniqueness_FirstCollision(t *testing.T) {
//	p := message.NewPrinter(language.English)
//
//	for i := 0; i < 25; i++ {
//		collisionAt := testFirstCollision(8)
//		t.Logf("first collision occurred after %s IDs.", p.Sprintf("%d", collisionAt))
//	}
//
//	t.Fatalf("test ended")
//}
//
//func testFirstCollision(idLen int) int {
//	seen := make(map[string]struct{})
//
//	count := 0
//	for {
//		id := GenerateID(idLen)
//		count++
//
//		if _, ok := seen[id]; ok {
//			return count // collision occurred here
//		}
//
//		seen[id] = struct{}{}
//	}
//}

func TestGenerateID_ZeroLength(t *testing.T) {
	id := GenerateID(0)
	if id != "" {
		t.Fatalf("expected empty string, got %q", id)
	}
}

func contains(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

package middlewares

import "testing"

func TestRemoveDiplicatesInURI(t *testing.T) {
	text := "123456//788888999990////"
	wants := "123456/788888999990/"
	if removeDuplicateSlashInURI(text) != wants {
		t.Errorf("Product of removeDuplicatesInURI is %s, wants %s", removeDuplicateSlashInURI(text), wants)
	}
}

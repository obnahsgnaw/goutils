package hsutil

import "testing"

func TestHash(t *testing.T) {
	raw := []byte("hello world")

	hashed, err := Md5(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("md5", string(hashed))

	hashed, err = Sha1(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("sha1", string(hashed))

	hashed, err = Sha256(raw)
	if err != nil {
		t.Fatal(err)
		return
	}
	println("sha256", string(hashed))
}

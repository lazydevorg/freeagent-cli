package cache

import "testing"

func TestCache(t *testing.T) {
	t.Setenv("CACHE_PATH", t.TempDir())

	err := Save("test", []byte("testdata"))
	if err != nil {
		t.Error(err)
	}
	data, err := Load("test")
	if err != nil {
		t.Error(err)
	}
	if string(data) != "testdata" {
		t.Errorf("Loaded data do not match: %s", string(data))
	}
}

func BenchmarkCache(b *testing.B) {
	b.Setenv("CACHE_PATH", b.TempDir())

	for i := 0; i < b.N; i++ {
		err := Save("test", []byte("testdata"))
		if err != nil {
			b.Error(err)
		}
		_, err = Load("test")
		if err != nil {
			b.Error(err)
		}
	}
}

func TestCacheJson(t *testing.T) {
	t.Setenv("CACHE_PATH", t.TempDir())

	type TestData struct {
		Name string
	}
	data := TestData{Name: "test"}
	err := SaveJson("test", data)
	if err != nil {
		t.Error(err)
	}
	var loadedData TestData
	err = LoadJson("test", &loadedData)
	if err != nil {
		t.Error(err)
	}
	if loadedData != data {
		t.Error("Loaded data do not match")
	}
}

func BenchmarkJsonCache(b *testing.B) {
	b.Setenv("CACHE_PATH", b.TempDir())

	type TestData struct {
		Name string
	}
	data := TestData{Name: "test"}

	for i := 0; i < b.N; i++ {
		err := SaveJson("test", data)
		if err != nil {
			b.Error(err)
		}
		var loadedData TestData
		err = LoadJson("test", &loadedData)
		if err != nil {
			b.Error(err)
		}
	}
}

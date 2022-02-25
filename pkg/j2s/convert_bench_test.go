package j2s_test

import (
	"testing"

	"github.com/masakurapa/go-json2struct/pkg/j2s"
)

func BenchmarkConvert(b *testing.B) {
	input := `{
		"key1": "aaaaa",
		"key2": 12345,
		"key3": 12345.6789,
		"key4": true,
		"key5": null,
		"key1": ["1", "2", "3", "4", "5", "6", "7", "8"]
	}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := j2s.Convert(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

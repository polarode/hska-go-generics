package quick

import (
	"bytes"
	"math/rand"
	"reflect"
)

func RandomStringGenerator(r *rand.Rand, size int, alphabet string) reflect.Value {
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := r.Intn(len(alphabet))
		buffer.WriteString(string(alphabet[index]))
	}
	return reflect.ValueOf(buffer.String())
}

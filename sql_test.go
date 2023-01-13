package sqlxcache

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_hashSQLQuery(t *testing.T) {
	type args struct {
		query  string
		params []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "example1", args: args{query: "select * from *", params: []any{"test1", 1}}, want: "c455f188767b0a753d7e0bda9f5294406eace3cb145de8f402a0c5cd44c4c1c5"},
		{name: "example2", args: args{query: "select * from Login WHERE Name=?", params: []any{"test1", 1}}, want: "e041c5e8f205cae18689ff4af6890a31d874401825dd055de631467e2b0b7d0e"},
		{name: "example3", args: args{query: "select * from Login WHERE Name=%s", params: []any{"test1", 1}}, want: "ae99637f601d55688c089b36081fe6709dbca45405a2e8911c02a469a8ef3980"},
		{name: "example4", args: args{query: "*", params: []any{1}}, want: "1c91a00e5c31062eaf7ad37757ff4b5543406112582dbcc9faa8b59024e23fb1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashSQLQuery(tt.args.query, tt.args.params...); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(got)
				t.Errorf("hashSQLQuery() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func Benchmark_hashSQLQuery(b *testing.B) {
	inputs := []struct {
		query string
		args  []any
	}{
		{query: "SELECT * FROM ?", args: []any{"test", 1}},
		{query: "SELECT Name,PW,Birthday FROM User WHERE Something= %d AND SomethingElse LIKE %s", args: []any{1, "test"}},
		{query: "SELECT" + strings.Repeat("Lorem ipsum ", 50) + "FROM Test WHERE x=? AND y=?", args: []any{"test", 1}},
		{query: "x", args: []any{"x", 1}},
	}

	for idx, v := range inputs {
		count := idx
		b.Run(fmt.Sprintf("hashSQLQuery_%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hashSQLQuery(v.query, v.args)
			}
		})
	}
}

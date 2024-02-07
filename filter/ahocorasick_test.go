package filter

import (
	"reflect"
	"testing"
)

var (
	words1 = []string{"敏感词1", "敏感词2", "敏感词3"}

	text1 = "敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词"
)

func Test_FindAll(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result []string
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: []string{"敏感词1", "敏感词2", "敏感词3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			matchAll := filter.FindAll(tt.args.text)
			if !reflect.DeepEqual(matchAll, tt.result) {
				t.Errorf("FindAll() = %v, want %v", matchAll, tt.result)
			}
		})
	}
}

func Test_FindAllCount(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result map[string]int
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: map[string]int{
				"敏感词1": 2,
				"敏感词2": 1,
				"敏感词3": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			matchAll := filter.FindAllCount(tt.args.text)
			if !reflect.DeepEqual(matchAll, tt.result) {
				t.Errorf("FindAllCount() = %v, want %v", matchAll, tt.result)
			}
		})
	}
}

func Test_FindOne(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: "敏感词1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			matchAll := filter.FindOne(tt.args.text)
			if !reflect.DeepEqual(matchAll, tt.result) {
				t.Errorf("FindOne() = %v, want %v", matchAll, tt.result)
			}
		})
	}
}

func Test_IsSensitive(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			matchAll := filter.IsSensitive(tt.args.text)
			if !reflect.DeepEqual(matchAll, tt.result) {
				t.Errorf("IsSensitive() = %v, want %v", matchAll, tt.result)
			}
		})
	}
}

func Test_Replace(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: "****,这是****,这是****,这是****,这里没有敏感词",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			result := filter.Replace(tt.args.text, rune('*'))
			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf("Remove() = %v, want %v", result, tt.result)
			}
		})
	}
}

func Test_Remove(t *testing.T) {
	type args struct {
		words []string
		text  string
	}

	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "1",
			args: args{
				words: words1,
				text:  text1,
			},
			result: ",这是,这是,这是,这里没有敏感词",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewAcModel()

			filter.AddWords(tt.args.words...)

			result := filter.Remove(tt.args.text)
			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf("Remove() = %v, want %v", result, tt.result)
			}
		})
	}
}

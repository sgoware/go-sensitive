package sensitive

import (
	"reflect"
	"testing"
)

func Test_NewFilter(t *testing.T) {
	type args struct {
		storeOption  StoreOption
		filterOption FilterOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "memory+dfa",
			args: args{
				storeOption: StoreOption{
					Type: StoreMemory,
					Dsn:  "",
				},
				filterOption: FilterOption{
					Type: FilterDfa,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterManager := NewFilter(tt.args.storeOption, tt.args.filterOption)

			err := filterManager.GetStore().AddWord("敏感词1", "敏感词2", "敏感词3")
			if err != nil {
				t.Errorf("add sensitive word failed, err: %v", err)
			}

			isSensitive := filterManager.GetFilter().IsSensitive("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			if !reflect.DeepEqual(isSensitive, true) {
				t.Errorf("IsSensitive() = %v, want %v", isSensitive, true)
			}

			filtered := filterManager.GetFilter().Filter("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			if !reflect.DeepEqual(filtered, "这是,这是,这是,这是,这里没有敏感词") {
				t.Errorf("IsSensitive() = %v, want %v", filtered, "这是,这是,这是,这是,这里没有敏感词")
			}

			replaced := filterManager.GetFilter().Replace("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词", '*')
			if !reflect.DeepEqual(replaced, "这是****,这是****,这是****,这是****,这里没有敏感词") {
				t.Errorf("IsSensitive() = %v, want %v", replaced, "这是****,这是****,这是****,这是****,这里没有敏感词")
			}

			matchedOne := filterManager.GetFilter().FindOne("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			if !reflect.DeepEqual(matchedOne, "敏感词1") {
				t.Errorf("IsSensitive() = %v, want %v", matchedOne, "敏感词1")
			}

			matchedAll := filterManager.GetFilter().FindAll("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			if !reflect.DeepEqual(matchedAll, []string{"敏感词1", "敏感词2", "敏感词3"}) {
				t.Errorf("IsSensitive() = %v, want %v", matchedAll, []string{"敏感词1", "敏感词2", "敏感词3"})
			}

			matchedMap := filterManager.GetFilter().FindAllCount("这是敏感词1,这是敏感词2,这是敏感词3,这是敏感词1,这里没有敏感词")
			if !reflect.DeepEqual(matchedMap, map[string]int{"敏感词1": 2, "敏感词2": 1, "敏感词3": 1}) {
				t.Errorf("IsSensitive() = %v, want %v", matchedMap, map[string]int{"敏感词1": 2, "敏感词2": 1, "敏感词3": 1})
			}
		})
	}
}

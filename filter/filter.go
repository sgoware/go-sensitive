package filter

type (
	Model interface {
		// Filter 过滤铭感词
		Filter(text string) string
		// Replace 和谐敏感词
		Replace(text string, repl rune) string
		// IsSensitive 是否有敏感词
		IsSensitive(text string) bool
		// FindOne 找到一个敏感词
		FindOne(text string) string
		// FindAll 找到所有敏感词
		FindAll(text string) []string
		// FindAllCount 找到所有敏感词及出现次数
		FindAllCount(text string) map[string]int
	}
)

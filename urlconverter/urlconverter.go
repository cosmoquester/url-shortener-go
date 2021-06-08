package urlconverter

// URLConverter 는 LongURL과 ShortURL을 기록하고 변환하고 삭제할 수 있는 interface입니다.
type URLConverter interface {
	GetShortURL(string) (string, bool)
	GetLongURL(string) (string, bool)
	PutURL(string) (string, bool)
	DelURL(string) bool
}

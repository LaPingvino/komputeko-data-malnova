package komputeko

type Terminaro []Entry
type Entry struct {
	Wordtype     string
	Translations []Translation
}
type Translation struct {
	Language string
	Words    []Word
}
type Word struct {
	Written string
	Sources []string
	//	Frequency float32
}

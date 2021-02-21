package services

type RepoInterface interface {
	GetFiles() []map[string]interface{}
	Push(filename, content string) (string,string,string)
	Del(filepath,sha string) string
}
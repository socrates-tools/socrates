package convention

type Tool struct {
	Namespace string
	Directory string
}

var Tools = []Tool{
	{Namespace: "review", Directory: "review"},
	{Namespace: "agora", Directory: "agora"},
}

func ToolDirectory(namespace string) string {
	for _, tool := range Tools {
		if tool.Namespace == namespace {
			return tool.Directory
		}
	}
	return ""
}

package main

type ActionsStruct struct {
	InsertionActions []string
	DeletionActions  []string
	CombinedActions  []string
}

type ActionArgs struct {
	Name           string
	InsertionCount int
	DeletionCount  int
	IS             string
	IES            string
	DS             string
	DES            string
	HeShe          string
}

var Actions = ActionsStruct{
	InsertionActions: []string{
		"{{.Name}} baked {{.InsertionCount}} muffin{{.IS}}",
		"{{.Name}} baked {{.InsertionCount}} cookie{{.IS}}",
		"{{.Name}} built a sandcastle with {{.InsertionCount}} tower{{.IS}}",
		"{{.Name}} watched {{.InsertionCount}} cat video{{.IS}}",
		"{{.Name}} planted {{.InsertionCount}} flower{{.IS}}",
		"{{.Name}} meditated {{.InsertionCount}} minute{{.IS}}",
		"{{.Name}} read {{.InsertionCount}} page{{.IS}} in a book",
	},
	DeletionActions: []string{
		"{{.Name}} ate {{.DeletionCount}} muffin{{.DS}}",
		"{{.Name}} ate {{.DeletionCount}} cookie{{.DS}}",
		"{{.Name}} kicked {{.DeletionCount}} traffic cone{{.DS}}",
		"{{.Name}} washed {{.DeletionCount}} dish{{.DES}}",
		"{{.Name}} washed {{.DeletionCount}} spoon{{.DS}}",
	},
	CombinedActions: []string{
		"{{.Name}} baked {{.InsertionCount}} muffin{{.IS}} and ate {{.DeletionCount}} muffin{{.DS}}",
		"{{.Name}} baked {{.InsertionCount}} cookie{{.IS}} and ate {{.DeletionCount}} cookie{{.DS}}",
		"{{.Name}} washed {{.InsertionCount}} shirt{{.IS}} and dried {{.DeletionCount}} shirt{{.DS}}",
		"{{.Name}} planted {{.InsertionCount}} flower{{.IS}} and {{.DeletionCount}} vegetable{{.DS}}",
		"{{.Name}} knitted {{.InsertionCount}} mitten{{.IS}} and {{.DeletionCount}} sock{{.DS}}",
		"{{.Name}} napped {{.InsertionCount}} minute{{.IS}} and walked around the house {{.DeletionCount}} minute{{.DS}}",
		"{{.Name}} read {{.InsertionCount}} page{{.IS}} in a book and " +
			"torn {{.DeletionCount}} page{{.DS}} out of the book because {{.HeShe}} didn't like them",
	},
}

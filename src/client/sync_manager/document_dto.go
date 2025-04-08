package sync_manager

import "editor-service/node/ot"

type DocumentConfig struct {
	docId    string
	authorId string
	version  int
	title    string
	document ot.Doc
}

func NewDocumentConfig(docId string, authorId string, version int, title string, document ot.Doc) DocumentConfig {
	return DocumentConfig{
		docId:    docId,
		authorId: authorId,
		version:  version,
		title:    title,
		document: document,
	}
}

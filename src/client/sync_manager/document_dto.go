package sync_manager

import "editor-service/node/ot"

type DocumentConfig struct {
	docId    string
	authorId string
	version  int
	title    string
	Document ot.Doc
}

func NewDocumentConfig(docId string, authorId string, version int, title string, document ot.Doc) DocumentConfig {
	return DocumentConfig{
		docId:    docId,
		authorId: authorId,
		version:  version,
		title:    title,
		Document: document,
	}
}

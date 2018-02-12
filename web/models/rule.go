package models

type Rule struct {
	baseModel
}

func NewRule() Rule {
	mdl := Rule{}
	mdl.collectionName = "rules"
	return mdl
}

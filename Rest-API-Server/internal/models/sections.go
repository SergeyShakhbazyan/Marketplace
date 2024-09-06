package models

import "github.com/gocql/gocql"

type Section struct {
	SectionID      gocql.UUID    `json:"sectionID"`
	SectionType    string        `json:"sectionType"`
	SectionHeading string        `json:"sectionHeading"`
	Content        []interface{} `json:"sectionContent"`
}

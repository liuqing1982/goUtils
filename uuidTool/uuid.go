package uuidTool

import uuid "github.com/satori/go.uuid"

func UuidCovery(uid string) (string, error) {
	d, err := uuid.FromString(uid)
	return d.String(), err
}

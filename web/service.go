package web

import (
	"fmt"
	"frankinstore/store"
)

type Service struct{}

// REVU: don't think we need to return the Service
func StartService(part int, db store.Store) (*Service, error) {
	return nil, fmt.Errorf("NewService not implemented!")
}

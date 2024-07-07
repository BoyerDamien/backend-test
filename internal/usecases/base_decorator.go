package usecases

import "github.com/japhy-tech/backend-test/internal/gateways"

type Base struct {
	datastore gateways.IDatastore
}

func (b *Base) Init(datastore gateways.IDatastore) {
	b.datastore = datastore
}

func (b Base) Datastore() gateways.IDatastore {
	return b.datastore
}

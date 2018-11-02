package controller

import (
	"github.com/jharrington22/cluster-readers/pkg/controller/clusterreader"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, clusterreader.Add)
}

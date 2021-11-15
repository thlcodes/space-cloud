package mgo

import (
	"context"
	"strings"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Create inserts a document (or multiple when op is "all") into the database
func (m *Mongo) Create(ctx context.Context, col string, req *model.CreateRequest) (int64, error) {
	// Create a collection object
	db := m.getClient().Database(m.dbName)
	cols, err := db.ListCollectionNames(ctx, utils.M{"name": utils.M{"$regex": "^" + strings.ReplaceAll(col, "_", "[-_]") + "$"}})
	if err != nil {
		return 0, err
	}
	if len(cols) > 0 {
		col = cols[0]
	}
	collection := db.Collection(col)

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo create", map[string]interface{}{"col": col, "doc": req.Document, "op": req.Operation})

	switch req.Operation {
	case utils.One:
		// Insert single document
		_, err := collection.InsertOne(ctx, req.Document)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		// Insert multiple documents
		objs, ok := req.Document.([]interface{})
		if !ok {
			return 0, utils.ErrInvalidParams
		}

		res, err := collection.InsertMany(ctx, objs)
		if err != nil {
			return 0, err
		}

		return int64(len(res.InsertedIDs)), nil

	default:
		return 0, utils.ErrInvalidParams
	}
}

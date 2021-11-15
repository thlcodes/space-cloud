package mgo

import (
	"context"
	"errors"
	"strings"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Delete removes the document(s) from the database which match the condition
func (m *Mongo) Delete(ctx context.Context, col string, req *model.DeleteRequest) (int64, error) {
	db := m.getClient().Database(m.dbName)
	cols, err := db.ListCollectionNames(ctx, utils.M{"name": utils.M{"$regex": "^" + strings.ReplaceAll(col, "_", "[-_]") + "$"}})
	if err != nil {
		return 0, err
	}
	if len(cols) > 0 {
		col = cols[0]
	}
	collection := db.Collection(col)
	req.Find = sanitizeWhereClause(ctx, col, req.Find)
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo delete", map[string]interface{}{"col": col, "find": req.Find, "op": req.Operation})

	switch req.Operation {
	case utils.One:
		_, err := collection.DeleteOne(ctx, req.Find)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		res, err := collection.DeleteMany(ctx, req.Find)
		if err != nil {
			return 0, err
		}

		return res.DeletedCount, nil

	default:
		return 0, errors.New("Invalid operation")
	}
}

// DeleteCollection removes a collection from database`
func (m *Mongo) DeleteCollection(ctx context.Context, col string) error {
	return m.getClient().Database(m.dbName).Collection(col, &options.CollectionOptions{}).Drop(ctx)
}

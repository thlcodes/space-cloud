package mgo

import (
	"context"
	"errors"
	"strings"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Aggregate performs a mongo db pipeline aggregation
func (m *Mongo) Aggregate(ctx context.Context, col string, req *model.AggregateRequest) (interface{}, error) {
	db := m.getClient().Database(m.dbName)
	cols, err := db.ListCollectionNames(ctx, utils.M{"name": utils.M{"$regex": "^" + strings.ReplaceAll(col, "_", "[-_]") + "$"}})
	if err != nil {
		return nil, err
	}
	if len(cols) > 0 {
		col = cols[0]
	}
	collection := db.Collection(col)

	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo agregate", map[string]interface{}{"col": col, "pipeline": req.Pipeline, "op": req.Operation})

	switch req.Operation {
	case utils.One:
		var result map[string]interface{}

		cur, err := collection.Aggregate(ctx, req.Pipeline)
		if err != nil {
			return nil, err
		}
		defer func() { _ = cur.Close(ctx) }()

		if !cur.Next(ctx) {
			return nil, errors.New("No result found")
		}

		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}

		return result, nil

	case utils.All:
		results := []interface{}{}

		cur, err := collection.Aggregate(ctx, req.Pipeline)
		defer func() { _ = cur.Close(ctx) }()
		if err != nil {
			return nil, err
		}

		for cur.Next(ctx) {
			var doc map[string]interface{}
			err := cur.Decode(&doc)
			if err != nil {
				return nil, err
			}

			results = append(results, doc)
		}

		if err := cur.Err(); err != nil {
			return nil, err
		}

		return results, nil

	default:
		return nil, utils.ErrInvalidParams
	}
}

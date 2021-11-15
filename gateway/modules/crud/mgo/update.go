package mgo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spaceuptech/helpers"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// Update updates the document(s) which match the condition provided.
func (m *Mongo) Update(ctx context.Context, col string, req *model.UpdateRequest) (int64, error) {
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
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Mongo update", map[string]interface{}{"col": col, "find": req.Find})

	switch req.Operation {
	case utils.One:
		_, err := collection.UpdateOne(ctx, req.Find, req.Update)
		if err != nil {
			return 0, err
		}

		return 1, nil

	case utils.All:
		res, err := collection.UpdateMany(ctx, req.Find, req.Update)
		if err != nil {
			return 0, err
		}

		return res.MatchedCount, nil

	case utils.Upsert:
		doUpsert := true
		res, err := collection.UpdateOne(ctx, req.Find, req.Update, &options.UpdateOptions{Upsert: &doUpsert})
		if err != nil {
			return 0, err
		}

		return res.MatchedCount + res.UpsertedCount, nil

	default:
		return 0, utils.ErrInvalidParams
	}
}

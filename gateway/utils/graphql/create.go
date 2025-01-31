package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/spaceuptech/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func (graph *Module) generateWriteReq(ctx context.Context, field *ast.Field, token string, store map[string]interface{}) (model.RequestParams, []*model.AllRequest, []interface{}, error) {
	dbAlias, err := graph.GetDBAlias(ctx, field, token, store)
	if err != nil {
		return model.RequestParams{}, nil, nil, err
	}

	col := strings.TrimPrefix(field.Name.Value, "insert_")

	docs, err := extractDocs(ctx, field.Arguments, store)
	if err != nil {
		return model.RequestParams{}, nil, nil, err
	}
	reqParams, reqs, returningDocs, err := graph.processNestedFields(ctx, docs, dbAlias, col, token)
	if err != nil {
		return model.RequestParams{}, nil, nil, err
	}

	return reqParams, reqs, returningDocs, nil
}

func (graph *Module) prepareDocs(doc map[string]interface{}, schemaFields model.Fields) {
	// FieldIDs is the array of fields for which an unique id needs to be generated. These will only be done for those
	// fields which have the type ID.
	// FieldDates is the array of fields for which the current time needs to be set.
	fieldIDs := make([]string, 0)
	fieldDates := make([]string, 0)
	fieldDefaults := make(map[string]interface{})

	for fieldName, fieldSchema := range schemaFields {
		// Only process ID fields which are required
		if fieldSchema.Kind == model.TypeID && fieldSchema.IsFieldTypeRequired {
			fieldIDs = append(fieldIDs, fieldName)
		}

		if fieldSchema.IsCreatedAt || fieldSchema.IsUpdatedAt {
			fieldDates = append(fieldDates, fieldName)
		}

		if fieldSchema.IsDefault {
			defaultStringValue, isString := fieldSchema.Default.(string)
			if fieldSchema.Kind == model.TypeJSON && isString {
				var v interface{}
				_ = json.Unmarshal([]byte(defaultStringValue), &v)
				fieldDefaults[fieldName] = v
			} else {
				fieldDefaults[fieldName] = fieldSchema.Default
			}
		}
	}

	// Set the default values if the field isn't set already. This always need to happen first
	for field, defaultValue := range fieldDefaults {
		if _, p := doc[field]; !p {
			doc[field] = defaultValue
		}
	}

	// Set a new id for all those field ids which do not have the field set already
	for _, field := range fieldIDs {
		if _, p := doc[field]; !p {
			doc[field] = primitive.NewObjectID()
		}
	}

	// Set the current time for all fieldDates
	for _, field := range fieldDates {
		doc[field] = time.Now().UTC()
	}
}

func copyDoc(doc map[string]interface{}) map[string]interface{} {
	newDoc := make(map[string]interface{}, len(doc))
	for k, v := range doc {
		newDoc[k] = v
	}

	return newDoc
}

func (graph *Module) processNestedFields(ctx context.Context, docs []interface{}, dbAlias, col, token string) (model.RequestParams, []*model.AllRequest, []interface{}, error) {
	createRequests := make([]*model.AllRequest, 0)
	afterRequests := make([]*model.AllRequest, 0)
	var err error
	var reqParams model.RequestParams
	// Check if we can the schema for this collection
	schemaFields, p := graph.schema.GetSchema(dbAlias, col)
	if !p {
		// Return the docs as is if no schema is available
		return model.RequestParams{}, []*model.AllRequest{{Type: string(model.Create), Col: col, Operation: utils.All, Document: docs, DBAlias: dbAlias}}, docs, nil
	}

	returningDocs := make([]interface{}, len(docs))

	for i, docTemp := range docs {

		// Each document is actually an object
		doc := docTemp.(map[string]interface{})
		graph.prepareDocs(doc, schemaFields)

		newDoc := copyDoc(doc)

		// Iterate over each field of the document to see if has any linked fields that are present
		for fieldName, fieldValue := range doc {
			fieldSchema, p := schemaFields[fieldName]
			if !p || !fieldSchema.IsLinked {
				// Simply ignore if the field does not have a corresponding schemaFields or it isn't linked
				continue
			}

			// We are here means that the field is actually a linked value

			fromFieldSchema, p := schemaFields[fieldSchema.LinkedTable.From]
			if !p {
				// Ignore if the `from` key is not present in the schema
				continue
			}

			// Ignore if the field key is present in the linked table config. We don't support such operations.
			if fieldSchema.LinkedTable.Field != "" {
				continue
			}

			if fromFieldSchema.IsPrimary {
				// Ignore if the from field doesn't exist in the document
				if _, p := doc[fromFieldSchema.FieldName]; !p {
					continue
				}
			}

			// Lets populate an array of linked docs
			var linkedDocs []interface{}
			if !fieldSchema.IsList {
				temp, ok := fieldValue.(map[string]interface{})
				if !ok {
					return model.RequestParams{}, nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid format provided for linked field %s - wanted object got array", fieldName), nil, nil)
				}

				linkedDocs = []interface{}{temp}
			} else {
				temp, ok := fieldValue.([]interface{})
				if !ok {
					return model.RequestParams{}, nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("invalid format provided for linked field %s - wanted array got object", fieldName), nil, nil)
				}
				linkedDocs = temp
			}

			// Iterate over each linked doc
			for _, linkedDocTemp := range linkedDocs {
				// Each document is actually an object
				linkedDoc := linkedDocTemp.(map[string]interface{})

				linkedSchemaFields, p := graph.schema.GetSchema(fieldSchema.LinkedTable.DBType, fieldSchema.LinkedTable.Table)
				if !p {
					return model.RequestParams{}, nil, nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("schema not provided for table (%s). Check the link directive for field (%s) in table (%s)", fieldSchema.LinkedTable.Table, fieldSchema.FieldName, col), nil, nil)
				}

				graph.prepareDocs(linkedDoc, linkedSchemaFields)

				// Check if the `from` field is a primary key. If it is, we need to set that value in the `to` field
				// of the nested value. If it is not a primary key, we'll have to set it with the value of the `to`
				// field of the nested value
				if fromFieldSchema.IsPrimary {
					linkedDoc[fieldSchema.LinkedTable.To] = doc[fieldSchema.LinkedTable.From]
				} else {
					// The nested docs need to be inserted first in this case
					doc[fieldSchema.LinkedTable.From] = linkedDoc[fieldSchema.LinkedTable.To]
				}
			}

			_, linkedCreateRequests, returningLinkedDocs, err := graph.processNestedFields(ctx, linkedDocs, fieldSchema.LinkedTable.DBType, fieldSchema.LinkedTable.Table, token)
			if err != nil {
				return model.RequestParams{}, nil, nil, err
			}

			if fromFieldSchema.IsPrimary {
				// It the from field is primary, it means that the nested docs need to be inserted after the parent docs have been inserted
				afterRequests = append(afterRequests, linkedCreateRequests...)
			} else {
				// If the from field is not primary, it means that the nested docs need to be inserted before the parent docs
				createRequests = append(createRequests, linkedCreateRequests...)
			}

			// Delete the nested field. The schema module would throw an error otherwise
			delete(doc, fieldName)

			newDoc[fieldName] = returningLinkedDocs
		}
		r := &model.CreateRequest{Document: []interface{}{doc}, Operation: utils.All}
		reqParams, err = graph.auth.IsCreateOpAuthorised(ctx, graph.project, dbAlias, col, token, r)
		if err != nil {
			return model.RequestParams{}, nil, nil, err
		}
		// Store the mutated doc because of security rules in returning docs
		for key, value := range doc {
			newDoc[key] = value
		}

		returningDocs[i] = newDoc
	}
	createRequests = append(createRequests, &model.AllRequest{Type: string(model.Create), Col: col, Operation: utils.All, Document: docs, DBAlias: dbAlias})
	return reqParams, append(createRequests, afterRequests...), returningDocs, nil
}

func extractDocs(ctx context.Context, args []*ast.Argument, store utils.M) ([]interface{}, error) {
	for _, v := range args {
		switch v.Name.Value {
		case "docs":
			temp, err := utils.ParseGraphqlValue(v.Value, store)
			if err != nil {
				return nil, err
			}
			arr, ok := temp.([]interface{})
			if !ok {
				return nil, errors.New("docs should be of type array []")
			}
			return arr, nil
		}
	}

	return []interface{}{}, nil
}

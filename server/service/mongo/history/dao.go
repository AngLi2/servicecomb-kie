/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package history

import (
	"context"
	"github.com/apache/servicecomb-kie/pkg/model"
	"github.com/apache/servicecomb-kie/server/service"
	"github.com/apache/servicecomb-kie/server/service/mongo/session"
	"github.com/go-mesh/openlogging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getHistoryByKeyID(ctx context.Context, filter bson.M) ([]*model.KVDoc, error) {
	collection := session.GetDB().Collection(session.CollectionKVRevision)
	cur, err := collection.Find(ctx, filter, options.Find().SetSort(map[string]interface{}{
		"revision": -1,
	}))
	if err != nil {
		return nil, err
	}
	kvs := make([]*model.KVDoc, 0)
	var exist bool
	for cur.Next(ctx) {
		var elem model.KVDoc
		err := cur.Decode(&elem)
		if err != nil {
			openlogging.Error("decode error: " + err.Error())
			return nil, err
		}
		exist = true
		kvs = append(kvs, &elem)
	}
	if !exist {
		return nil, service.ErrRevisionNotExist
	}
	return kvs, nil
}

//AddHistory add kv history
func AddHistory(ctx context.Context, kv *model.KVDoc) error {
	ctx, cancel := context.WithTimeout(ctx, session.Timeout)
	defer cancel()
	collection := session.GetDB().Collection(session.CollectionKVRevision)
	_, err := collection.InsertOne(ctx, kv)
	if err != nil {
		openlogging.Error(err.Error())
		return err
	}
	return nil
}

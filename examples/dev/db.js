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

db.createUser(
    {
        user: "kie",
        pwd: "123",
        roles:[
            {
                role: "readWrite",
                db:   "kie"
            }
        ]
    }
);
db.createCollection("counter");
db.counter.insertOne( { name:"revision_counter",count: 1 } );
db.createCollection( "kv", {
    validator: { $jsonSchema: {
            bsonType: "object",
            required: [ "key","domain","project","id","value","create_revision","update_revision","value_type","label_id" ],
            properties: {
                key: {
                    bsonType: "string",
                },
                value_type: {
                    enum: [ "text","string","yaml", "json", "properties", "ini" ]
                },
                labels: {
                    bsonType: "object"
                }
            }
        } }
} );
db.kv.createIndex({"id": 1}, { unique: true } );
db.kv.createIndex({key: 1, label_id: 1,domain:1,project:1},{ unique: true });
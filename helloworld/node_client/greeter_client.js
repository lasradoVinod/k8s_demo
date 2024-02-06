/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

var PROTO_PATH = './helloworld.proto';

var fs = require('fs');
var parseArgs = require('minimist');
var grpc = require('@grpc/grpc-js');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
var hello_proto = grpc.loadPackageDefinition(packageDefinition).helloworld;

function sleep(milliseconds) {
  const start = new Date().getTime();
  let now = start;
  while (now - start < milliseconds) {
    now = Date.now();
  }
}

function main() {
  var argv = parseArgs(process.argv.slice(2))
  if (argv.target) {
    target = argv.target;
  } else {
    target = 'localhost:50051';
  }

  if (!argv.caKey || !argv.Cert || !argv.Key) {
    console.log("No values provided")
    process.exit(1);
  }

  const credentials = grpc.credentials.createSsl(
    fs.readFileSync(argv.caKey), 
    fs.readFileSync(argv.Key), 
    fs.readFileSync(argv.Cert)
  );

  while (true) {
    console.log (hello_proto)
    var client = new hello_proto.Greeter(target, credentials);
    console.log(client)
    for (let i = 0; i <100; i++) {
      console.log(i.toString())
      client.sayHello({name: i.toString()}, (err, response) => {
        if (err) {
          console.log(err)
        } else {
          console.log('Greeting:', response.message);
        }
      });
      sleep(100);
    }
  }
}

main();

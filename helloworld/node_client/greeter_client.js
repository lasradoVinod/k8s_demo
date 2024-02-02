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
  while (new Date().getTime() - start < milliseconds);
}


function main() {
  var argv = parseArgs(process.argv.slice(2))
  if (argv.target) {
    target = argv.target;
  } else {
    target = 'localhost:50051';
  }

  if (!argv.caKey || !argv.sCert || !argv.sKey) {
    process.exit(1);
  }

  const credentials = grpc.credentials.createSsl(
    fs.readFileSync(argv.caKey), 
    fs.readFileSync(argv.sKey), 
    fs.readFileSync(argv.sCert)
  );

  while (true) {
    var client = new hello_proto.Greeter(target, credentials);
    var user;
    if (argv._.length > 0) {
      user = argv._[0];
    } else {
      user = 'world';
    }
    for (let i = 0; i <100; i++) {
      client.sayHello({name: user}, function(err, response) {
        console.log('Greeting:', response.message);
      });
      sleep(100);
    }
  }
}

main();

// Copyright 2020 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START cloudrun_pubsub_server_setup]
// [START run_pubsub_server_setup]
const express = require('express');
const {VM} = require('vm2');
var util = require('util');
const axios = require('axios');
const axiosRetry = require('axios-retry');

const app = express();

axiosRetry(axios, { retryDelay: axiosRetry.exponentialDelay});

MAX_LEN = 100
API_HOST = process.env?.API_HOST
API_KEY = process.env?.API_KEY

const config = {
  headers: { Authorization: `Bearer ${API_KEY}` }
};

function LoggingFromVM(type, args, logs) {
  if (logs.length === MAX_LEN) {
    logs.push("...")
    return
  }
  if (logs.length >MAX_LEN) {
    return
  }
  logs.push(`[${type}] ${args.map((x) => util.inspect(x)).join(' ')}`)
}

async function RunUserCode(code) {
  const logs = [];
  const vm = new VM({
    timeout: 10000,
    sandbox: {
      console: {
          log: (...args) => LoggingFromVM('LOG', args, logs),
          warn: (...args) => LoggingFromVM('WARN', args, logs),
          error: (...args) => LoggingFromVM('ERROR', args, logs),        
      }
    }
  });

  try {
    vm.run(code);
  } catch (err) {
    console.log(err)
    return 'Failed to execute script.'
  }

  return logs.join('\n')
}

async function SendResult(result, sessionId, codeId) {
  try {
    await axios.post(
      `${API_HOST}/conference-session/${sessionId}/code/${codeId}/`,
      {result},
      config
    )
  } catch (err) {
    console.error(err)
  }
}


// This middleware is available in Express v4.16.0 onwards
app.use(express.json());
// [END run_pubsub_server_setup]
// [END cloudrun_pubsub_server_setup]

// [START cloudrun_pubsub_handler]
// [START run_pubsub_handler]
app.post('/', async(req, res) => {
  if (!req.body) {
    const msg = 'no Pub/Sub message received';
    console.error(`error: ${msg}`);
    res.status(400).send(`Bad Request: ${msg}`);
    return;
  }
  if (!req.body.message) {
    const msg = 'invalid Pub/Sub message format';
    console.error(`error: ${msg}`);
    res.status(400).send(`Bad Request: ${msg}`);
    return;
  }
  const bufferedData = req.body.message.data
  const dataString = Buffer.from(bufferedData, 'base64').toString().trim()
  const data = JSON.parse(dataString)
  const result = await RunUserCode(data.code)
  console.log('Result', result)
  await SendResult(result, data.sessionId, data.id)

  res.status(204).send();
});
// [END run_pubsub_handler]
// [END cloudrun_pubsub_handler]

module.exports = app;
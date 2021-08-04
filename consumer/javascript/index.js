// Imports the Google Cloud client library
const {PubSub} = require('@google-cloud/pubsub');
const {VM} = require('vm2');
var util = require('util');
const axios = require('axios');
const axiosRetry = require('axios-retry');

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

async function start({
  projectId, // Your Google Cloud Platform project ID
  topicName, // Name for the new topic to create
  subscriptionName, // Name for the new subscription to create
}) {
  // Instantiates a client
  const pubsub = new PubSub({projectId});

  // Creates a new topic
  var topic
  try {
    [topic] = await pubsub.createTopic(topicName);
  } catch (err) {
    topic = pubsub.topic(topicName);
  }
  console.log(`Topic ${topic.name} created.`);

  var subscription
  // Creates a subscription on that new topic
  try{
    [subscription] = await topic.createSubscription(subscriptionName);
  } catch (err) {
    subscription = topic.subscription(subscriptionName)
  }

  // Receive callbacks for new messages on the subscription
  subscription.on('message', async(message) => {
    const jsString = message.data.toString()
    const data = JSON.parse(jsString)
    const result = await RunUserCode(data.code)
    console.log('Result', result)
    await SendResult(result, data.sessionId, data.id)
    message.ack();
  });

  // Receive callbacks for errors on the subscription
  subscription.on('error', error => {
    console.error('Received error:', error);
  });
}

(async ()=>{
  await start({
    projectId: process.env?.GOOGLE_CLOUD_PROJECT,
    topicName: 'javascript_topic',
    subscriptionName: 'javascript_topic_consumer',
  });
})();

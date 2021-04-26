// Imports the Google Cloud client library
const {PubSub} = require('@google-cloud/pubsub');

async function start({
  projectId, // Your Google Cloud Platform project ID
  topicName, // Name for the new topic to create
  subscriptionName, // Name for the new subscription to create
}) {
  // Instantiates a client
  const pubsub = new PubSub({projectId});

  // Creates a new topic
  const [topic] = await pubsub.createTopic(topicName);
  console.log(`Topic ${topic.name} created.`);

  // Creates a subscription on that new topic
  const [subscription] = await topic.createSubscription(subscriptionName);

  // Receive callbacks for new messages on the subscription
  subscription.on('message', message => {
    console.log('Received message:', message.data.toString());
  });

  // Receive callbacks for errors on the subscription
  subscription.on('error', error => {
    console.error('Received error:', error);
  });
}

(async ()=>{
  await start({
    projectId: process.env?.PUBSUB_PROJECT_ID,
    topicName: 'javascript_topic',
    subscriptionName: 'javascript_topic',
  });
})();

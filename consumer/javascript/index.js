// Imports the Google Cloud client library
const {PubSub} = require('@google-cloud/pubsub');
const app = require('./app.js');
const PORT = process.env.PORT || 8081;


const createTopicSubscription = async ()=>{
  // Creates a client
  const pubsub = new PubSub();
  const subscriptionName = "js-consumer"
  const topicName = process.env.JAVASCRIPT_TOPIC

  var pushEndpoint = (process.env.BACKEND_TYPE === 'dev')? `http://javascriptconsumer:${PORT}`: `js-consumer.rfrl.ca`

  const options = {
    pushConfig: {
      // Set to your local endpoint.
      pushEndpoint,
    },
  };
  
  // Creates a new topic
  var topic
  try {
    [topic] = await pubsub.createTopic(topicName);
    console.log(`Topic ${topic.name} created.`);
  } catch (err) {
    topic = pubsub.topic(topicName);
  }
  
  try {
    await topic.createSubscription(subscriptionName, options);
  } catch (err) {
    console.error(err)
    return
  }
}

createTopicSubscription().catch((err)=> {console.log(err)});

app.listen(PORT, () =>
  console.log(`js-consumer listening on port ${PORT}`)
);
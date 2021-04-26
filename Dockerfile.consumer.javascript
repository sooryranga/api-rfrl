FROM node:15

# Create app directory
WORKDIR /app

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available (npm@5+)
COPY consumer/javascript/package*.json ./

RUN npm ci --only=production
# If you are building your code for production
# RUN npm ci --only=production

RUN npm install -g nodemon

CMD [ "nodemon", "consumer.js" ]
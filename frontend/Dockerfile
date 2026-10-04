FROM node:24.18.0 AS build

WORKDIR /app
COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000", "--mode", "devnetwork"]
FROM node:18-alpine

WORKDIR /var/app
COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN rm -rf .next
RUN npm run build

CMD ["npm", "start"]
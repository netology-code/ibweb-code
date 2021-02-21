FROM node:lts
WORKDIR /app
COPY package*.json /app/
RUN npm ci
COPY . /app
CMD ["npm", "start"]

# Это всего лишь документация
EXPOSE 9999

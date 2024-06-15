FROM node:20.14.0-slim

WORKDIR /src/
COPY internal/client ./client
COPY internal/docs ./docs

WORKDIR /src/client
RUN npm install --include=dev

CMD ["node", "--version"]

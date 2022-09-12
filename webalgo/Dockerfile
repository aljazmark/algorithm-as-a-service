FROM node:13.12.0-alpine as build
RUN apk add --no-cache bash
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY package.json ./
COPY package-lock.json ./
COPY default.conf.template ./
RUN npm ci --silent
RUN npm install react-scripts@4.0.1 -g --silent
COPY . ./
RUN npm run build

# production environment
FROM nginx:stable-alpine
RUN apk add --no-cache bash
COPY --from=build /app/build /usr/share/nginx/html
COPY --from=build /app/default.conf.template /etc/nginx/conf.d/default.conf.template
EXPOSE 3000
CMD /bin/bash -c "envsubst '\$PORT' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf" && nginx -g 'daemon off;'
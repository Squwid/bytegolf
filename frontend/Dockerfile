FROM node:18.12.1-alpine3.15 as build

COPY . .

# TODO: Remove this line after updating deps
ENV NODE_OPTIONS --openssl-legacy-provider 
RUN npm install
RUN npm run build:master

FROM nginx:alpine
COPY --from=build nginx.conf /etc/nginx/conf.d/default.conf

RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /build/ /usr/share/nginx/html

ENV PORT 8080
EXPOSE ${PORT}

CMD ["nginx", "-g", "daemon off;"]
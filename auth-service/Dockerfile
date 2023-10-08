FROM node:latest as build
WORKDIR /app
ADD . .
RUN npm install --omit-dev

FROM alpine:latest as main
COPY --from=build /app /
EXPOSE 3000
CMD ["npm","start"]

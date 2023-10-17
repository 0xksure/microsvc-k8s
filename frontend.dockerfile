FROM node:18-alpine as  build

WORKDIR /app
COPY . .
RUN yarn install
RUN yarn run build

FROM node:18-alpine as production
COPY --from=build /app/build .
COPY --from=build /app/package.json .
COPY --from=build /app/node_modules ./node_modules 
EXPOSE 3000
CMD ["node", "."]

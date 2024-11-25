FROM node:20.8.1 AS base
WORKDIR /base
COPY package.json ./
COPY package-lock.json ./
COPY . .

RUN npm install


FROM base AS build
ENV NODE_ENV=production
WORKDIR /build
COPY --from=base /base ./
RUN npm run build

# STAGE 3: app
FROM node:20.8.1 As app
ENV NODE_ENV=production
WORKDIR /app
COPY --from=build /build/.next ./.next
COPY --from=build /build/public ./public
COPY --from=build /build/node_modules ./node_modules
COPY --from=build /build/package.json ./

EXPOSE 3000

CMD ["npm", "start"]


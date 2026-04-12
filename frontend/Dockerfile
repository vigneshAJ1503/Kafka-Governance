FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
# Since we don't have package.json dependencies lock generated originally in this sandbox,
# we use npm install to generate everything from scratch smoothly.
RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV production

COPY --from=builder /app/next.config.js ./

COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000
ENV PORT 3000

CMD ["node", "server.js"]

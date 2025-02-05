FROM golang:latest AS builder
WORKDIR /app
COPY backend/ ./
RUN go build -o app .

FROM node:20 AS react-builder
WORKDIR /frontend
COPY frontend/ ./
RUN npm install
RUN npm run build

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=react-builder /frontend/dist ./frontend/dist
EXPOSE 8081

CMD ["./app"]

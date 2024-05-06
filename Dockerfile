FROM alpine
RUN apk add entr
RUN mkdir /app
WORKDIR /app
COPY bin/out/rinha-de-backend-2024-q1 ./
CMD ls rinha-de-backend-2024-q1 | entr -rn ./rinha-de-backend-2024-q1

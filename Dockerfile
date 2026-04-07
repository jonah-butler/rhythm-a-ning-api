FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY main .

RUN chmod +x ./main

EXPOSE 8080

USER appuser

CMD ["./main"]
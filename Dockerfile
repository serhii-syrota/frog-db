FROM go:latest
RUN go build .
CMD [ "./frog.db" ]
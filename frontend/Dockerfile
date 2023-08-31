# Use an official Go image as the base image
FROM golang:1.19

# Set the working directory
WORKDIR /app

# Copy the source code
ARG CACHEBUST
COPY . .
RUN echo $CACHEBUST

# Build the frontend binary and WebAssembly
RUN GOARCH=wasm GOOS=js go build -o web/app.wasm
RUN go build

# Run the built frontend binary (this will not be the final use)
CMD ["./frontend"]

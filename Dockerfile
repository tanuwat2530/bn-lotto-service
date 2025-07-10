# ---- Build Stage ----
# Use the official Go image as a builder.
# Using a specific version is recommended for reproducibility.
FROM golang:1.22.1-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
# This leverages Docker's layer caching. Dependencies are only re-downloaded
# if these files change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the Go application.
# -o /app/server specifies the output file name.
# CGO_ENABLED=0 disables Cgo, which is needed for a static binary.
# GOOS=linux specifies the target operating system.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/bn-lotto-app .

# ---- Final Stage ----
# Use a minimal, non-root base image for the final container.
# "distroless" images contain only your application and its runtime dependencies.
# They do not contain package managers, shells, or other programs.
FROM gcr.io/distroless/static-debian12

# Set the working directory.
WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/bn-lotto-app .

# Set the user to a non-root user for security.
# The `nonroot` user and group (UID/GID 65532) are provided by the distroless image.
USER nonroot:nonroot

# Expose the port the app runs on. Cloud Run will automatically
# use the PORT environment variable.
EXPOSE 8080

# Command to run the executable.
ENTRYPOINT ["/app/bn-lotto-app"]

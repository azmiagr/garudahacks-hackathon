#!/bin/bash
set -e

IMAGE="ghcr.io/${GITHUB_REPOSITORY:-azmiagr/garudahacks-hackathon}:latest"

echo "Syncing latest config from repository..."
git pull origin main

echo "Pulling latest image: $IMAGE"
docker pull "$IMAGE"

echo "Restarting services..."
docker compose down --remove-orphans
docker compose up -d

echo "Cleaning up unused images..."
docker image prune -f

echo "Deploy complete."
docker compose ps

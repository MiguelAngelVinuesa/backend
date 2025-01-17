#!/bin/bash
killall kubectl
kubectl port-forward deployment/d-store 8001:8080 &
kubectl port-forward deployment/i18n-svc 8002:8080 &
kubectl port-forward deployment/game-graph 8003:8080 &
kubectl port-forward deployment/game-service 8004:8080 &

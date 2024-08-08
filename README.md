# Command:
kubectl apply -f bookserver-deployment.yaml

kubectl port-forward svc/book-service 8080:3200


Now go to http://localhost:8080/books

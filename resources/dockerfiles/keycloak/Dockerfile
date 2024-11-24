FROM quay.io/keycloak/keycloak:26.0
ADD ./ /opt/keycloak/data/import
EXPOSE 8080
CMD ["start-dev", "--import-realm"]

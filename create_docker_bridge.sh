docker network create multi-host-network

docker network connect multi-host-network login-service

docker network connect multi-host-network postgres12


docker network connect multi-host-network playlist-postgres

docker network connect multi-host-network playlist-service

docker network connect multi-host-network subscription-service

docker network connect multi-host-network subscription-postgres

docker network connect multi-host-network mail-service-mailhog-1 

docker network connect multi-host-network mail-service-mailer-service-1 

docker network rm multi-host-network
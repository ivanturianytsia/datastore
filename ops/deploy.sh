export DIGITALOCEAN_ACCESS_TOKEN=
docker-machine create \
  --driver digitalocean \
  --digitalocean-size 512mb \
  --digitalocean-ssh-key-fingerprint 1b:2d:18:5f:96:5d:3d:7b:ea:c5:40:6f:69:bc:cb:b4 \
  agh-datastore

eval $(docker-machine env agh-datastore)

docker swarm init --advertise-addr=$(docker-machine ip agh-datastore)

export MAILGUN_DOMAIN=
export MAILGUN_KEY=
export MAILGUN_PUBKEY=

docker stack deploy --compose-file ./ops/docker-compose.yml agh

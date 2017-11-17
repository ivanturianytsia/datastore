export APPNAME="agh-datastore"
export REGISTRY="104.131.38.166:5000"
export HASH=$(git rev-parse HEAD)

# Deploy settings
export TRAEFIK_NETWORK="traefik-net"
export DOMAIN="datastore.agh.edu.pl"

# App-specific
export TOKEN_SECRET="THSIISaSecrettokprodenSecRet"
export TOKEN_EXPIRATION="72"
export ENV="prod"
export PORT="3000"
export STATIC_DIR="./client/dist"
export DB="mongodb://ivan:ivan@cluster0-shard-00-00-p1nri.mongodb.net:27017,cluster0-shard-00-01-p1nri.mongodb.net:27017,cluster0-shard-00-02-p1nri.mongodb.net:27017/agh?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"


function start {
  Y='\033[1;33m'
  NC='\033[0m'
  printf "${Y}\n - $@...\n${NC}"
}
function complete {
  G='\033[0;32m'
  NC='\033[0m'
  printf "${G}\n - Completed: $@.\n${NC}"
}

function build {
  case $1 in
    client)
      STEPNAME="Building client"
      COMMAND="cd client && \
        yarn run build && \
        cd .."
      ;;
    mac)
      STEPNAME="Building macOS binary"
      COMMAND="go get -v -d -t ./app && \
        go test --cover -v ./app/... && \
        go build -v -o bin/${APPNAME}_mac ./app"
      ;;
    alpine)
      STEPNAME="Building Alpine Linux binary"
      COMMAND="go get -v -d -t ./app && \
        go test --cover -v ./app/... && \
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/${APPNAME}_alpine ./app"
      ;;
    image)
      STEPNAME="Building Docker image"
      COMMAND="docker build \
        --tag ${APPNAME}:latest \
        --file ./ops/Dockerfile \
        --build-arg APPNAME=${APPNAME} \
        ."
      ;;
    push)
      STEPNAME="Push Docker image to Registry"
      COMMAND="docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:latest && \
        docker push $REGISTRY/${APPNAME}:latest && \
        docker tag ${APPNAME}:latest $REGISTRY/${APPNAME}:$HASH && \
        docker push $REGISTRY/${APPNAME}:$HASH"
      ;;
    service_update)
      STEPNAME="Update Docker service image"
      COMMAND="docker service update \
        --image $REGISTRY/${APPNAME}:$HASH \
        ${APPNAME}"
      ;;
    service_create_traefik)
      STEPNAME="Create Docker service"
      COMMAND="docker service create \
        --name ${APPNAME} \
        --label traefik.port=3000 \
        --label traefik.frontend.rule=Host:${DOMAIN} \
        --network ${TRAEFIK_NETWORK} \
        --env DB=$DB \
        --env DATA_DIR=$DATA_DIR \
        --env DIST_DIR=$DIST_DIR \
        --env CLIENT_ID=$CLIENT_ID \
        --env CLIENT_SECRET=$CLIENT_SECRET \
        --env REDIRECT_URI=$REDIRECT_URI \
        $REGISTRY/${APPNAME}:$HASH"
      ;;
    service_create)
      STEPNAME="Create Docker service"
      COMMAND="docker service create \
        --name ${APPNAME} \
        --publish 3000 \
        --network ${TRAEFIK_NETWORK} \
        --env DB=$DB
        $REGISTRY/${APPNAME}:$HASH && \
        docker service inspect ${APPNAME}"
      ;;
  esac

  if [ "$2" = "print" ]
  then
    echo $COMMAND
  else
    eval $COMMAND
    complete $STEPNAME
  fi
}

if [ "$1" = "" ]
then
  build alpine
  build image
  build push
else
  build $1 $2
fi

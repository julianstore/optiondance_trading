#!/bin/bash

PS3='Please enter your choice system to deploy: '
options=("option-dance-webapp" "option-dance-api" "option-dance-engine" "option-dance-engine-node1"   "Quit")
select opt in "${options[@]}"
do
    case $opt in
        "option-dance-webapp")
            echo "you chose option-dance-webapp,starting deploy..."
            cd /home/optiondance/option-dance/webapp
            git reset --hard
            git pull
            npm run build:beta
            ;;
        "option-dance-api")
            echo "you chose api ,starting deploy..."
            cd /home/optiondance/option-dance/server
            git reset --hard
            git pull
            go build -o optiondance
            systemctl restart option-dance-api
            ;;
        "option-dance-engine")
            echo "you chose engine master, starting deploy.."
            cd /home/optiondance/option-dance/server
            git reset --hard
            git pull
            go build -o optiondance
            systemctl restart option-dance-engine
            ;;
        "option-dance-engine-node1")
            echo "you chose engine node1, starting deploy.."
            cd /home/optiondance/option-dance/server
            git reset --hard
            git pull
            go build -o admin -i main.go
            go build -o optiondance
            systemctl restart option-dance-engine-node1
            ;;
        "Quit")
            break
            ;;
        *) echo "invalid option $REPLY";;
    esac
done

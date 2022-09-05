# Nexus

Запуск:
```bash
docker run -d -p 192.168.56.10:8081:8081 -p 192.168.56.10:8082:8082 --name nexus -e INSTALL4J_ADD_VM_PARAMS="-Xms512m -Xmx512m -XX:MaxDirectMemorySize=273m" sonatype/nexus3
```

Вывести пароль администратора:
```
docker exec -t nexus bash -c 'cat /nexus-data/admin.password && echo'
```

# Jenkins

Jenkins url: http://192.168.56.10:8080/
Repo url: https://github.com/kozl/netology-devops.git

## Install go

https://go.dev/doc/install

```
wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
```

## Freestyle job

1. `go test .`
2. `docker build . -t ubuntu-bionic:8082/hello-world:v$BUILD_NUMBER`
3. `docker login ubuntu-bionic:8082 -u admin -p admin && docker push ubuntu-bionic:8082/hello-world:v$BUILD_NUMBER && docker logout`

## Pipeline

```
pipeline {
 agent any
 stages {
  stage('Git') {
   steps {git 'https://github.com/dmitriy-tomin/netology-devops-cicd'}
  }
  stage('Test') {
   steps {
    sh 'go test .'
   }
  }
  stage('Build') {
   steps {
    sh 'docker build . -t ubuntu-bionic:8082/hello-world:v$BUILD_NUMBER'
   }
  }
  stage('Push') {
   steps {
    sh 'docker login ubuntu-bionic:8082 -u admin -p admin && docker push ubuntu-bionic:8082/hello-world:v$BUILD_NUMBER && docker logout'   }
  }
 }
}
```

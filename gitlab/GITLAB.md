# Gitlab

## Развернуть инсталляцию

1. Добавить запись в локальный /etc/hosts:

```bash
echo '192.168.56.10    gitlab.localdomain gitlab' >> /etc/hosts
```

2. Запустить инсталляцию:

```bash
VAGRANT_EXPERIMENTAL="disks" vagrant up
```
Переменная окружения VAGRANT_EXPERIMENTAL нужна для того, чтобы vagrant настроил виртуалку с нестандартным размером диска.

После успешного завершения Gitlab будет доступен по адресу http://gitlab.localdomain

Получить первичный пароль для пользователя root:

```bash
vagrant ssh -- sudo cat /etc/gitlab/initial_root_password
```

### Возможные проблемы

1. В Vagrantfile определено, что виртуальная машина, которая будет поднята — получит адрес 192.168.56.10. Возможна ситуация, когда сетевой интерфейс virtualbox, vboxnet0, использует другую подсеть:

```bash
$ ip -br -c a sh vboxnet2
vboxnet2         UP             192.168.100.1/24 fe80::800:27ff:fe00:2/64
```

В этом случае нужно исправить адрес на тот, что будет соответствовать сети интерфейса и добавить запись для него в /etc/hosts

2. Недостаточно ресурсов. Для разворачивания инсталляции потребуется как минимум 6GB памяти, следует учитывать это при запуске. Если ресурсов недостаточно — можно попробовать подправить Vagrantfile, но работоспособность инсталляции при этом не гарантируется. 

## Gitlab runner

Регистрация раннера:
```bash
   docker run -ti --rm --name gitlab-runner \
     --network host \
     -v /srv/gitlab-runner/config:/etc/gitlab-runner \
     -v /var/run/docker.sock:/var/run/docker.sock \
     gitlab/gitlab-runner:latest register
```

Конфигурация раннера для docker-in-docker:
```yaml
    volumes = ["/cache", "/var/run/docker.sock:/var/run/docker.sock"]
    extra_hosts = ["gitlab.localdomain:192.168.56.10"]
```

Запуск:
```bash
   docker run -d --name gitlab-runner --restart always \
     --network host \
     -v /srv/gitlab-runner/config:/etc/gitlab-runner \
     -v /var/run/docker.sock:/var/run/docker.sock \
     gitlab/gitlab-runner:latest
```

## Pipeline

```yaml
stages:
  - test
  - build

test:
  stage: test
  image: golang:1.17
  script: 
   - go test .

build:
  stage: build
  image: docker:latest
  script:
   - docker build .
```

## Sonarqube

```bash
vagrant ssh
cd /vagrant
docker-compose up -d
```

Sonarqube будет доступен по адресу http://gitlab.localdomain:9000

## Pipeline with sonarqube-check

```yaml
stages:
  - test
  - static-analysis
  - build

test:
  stage: test
  image: golang:1.16
  script: 
   - go test .

static-analysis:
 stage: test
 image:
  name: sonarsource/sonar-scanner-cli
  entrypoint: [""]
 variables:
 script:
  - sonar-scanner -Dsonar.projectKey=netology-gitlab -Dsonar.sources=. -Dsonar.host.url=http://gitlab.localdomain:9000 -Dsonar.login=a778675a32f0d9d6455a3d502f4e2838e784994d

build:
  stage: build
  image: docker:latest
  script:
   - docker build .
```

## Pipeline with manual run

```yaml
stages:
  - test
  - build

test:
  stage: test
  image: golang:1.16
  script: 
   - go test .

sonarqube-check:
 stage: test
 image:
  name: sonarsource/sonar-scanner-cli
  entrypoint: [""]
 variables:
 script:
  - sonar-scanner -Dsonar.projectKey=netology-gitlab -Dsonar.sources=. -Dsonar.host.url=http://gitlab.localdomain:9000 -Dsonar.login=a778675a32f0d9d6455a3d502f4e2838e784994d

build:
  stage: build
  image: docker:latest
  only:
    - master
  script:
   - docker build .

build:
  stage: build
  image: docker:latest
  when: manual
  except:
    - master
  script:
   - docker build .
```
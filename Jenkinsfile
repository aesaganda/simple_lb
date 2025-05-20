pipeline {
  agent {
    kubernetes {
      label 'docker-agent'
      defaultContainer 'jnlp'
      yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: jnlp
    image: jenkins/inbound-agent:latest
    args: ['\$(JENKINS_SECRET)', '\$(JENKINS_NAME)']
    tty: true
  - name: docker
    image: docker:24.0.5-cli
    command:
    - cat
    tty: true
    volumeMounts:
    - name: docker-sock
      mountPath: /var/run/docker.sock
  volumes:
  - name: docker-sock
    hostPath:
      path: /var/run/docker.sock
"""
    }
  }

  environment {
    DOCKER_IMAGE = "aesaganda/jenkinsci-rr:latest"
    DOCKER_CREDENTIALS_ID = 'dockerhub-creds' // Jenkins credentials ID for Docker Hub
  }

  stages {
    stage('Checkout') {
      steps {
        git url: 'https://github.com/aesaganda/simple_lb.git', branch: 'main'
      }
    }

    stage('Build Docker Image') {
      steps {
        container('docker') {
          script {
            sh """
              docker build -t ${DOCKER_IMAGE} -f go/rr/rr.Dockerfile .
            """
          }
        }
      }
    }

    stage('Push Docker Image') {
      steps {
        container('docker') {
          withCredentials([usernamePassword(credentialsId: "${DOCKER_CREDENTIALS_ID}", usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
            sh '''
              echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
              docker push ${DOCKER_IMAGE}
            '''
          }
        }
      }
    }
  }
}

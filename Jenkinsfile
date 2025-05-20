pipeline {
  agent {
    kubernetes {
      label 'docker-agent'
    }
  }

  environment {
    IMAGE_NAME = 'aesaganda/rrimage'
    IMAGE_TAG = 'latest'
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
          sh """
            docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f go/rr/rr.Dockerfile go/rr
          """
        }
      }
    }

    stage('Push Docker Image') {
      steps {
        container('docker') {
          withCredentials([usernamePassword(credentialsId: "${DOCKER_CREDENTIALS_ID}", usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
            sh '''
              echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin ${DOCKER_REGISTRY}
              docker push ${IMAGE_NAME}:${IMAGE_TAG}
            '''
          }
        }
      }
    }
  }
}

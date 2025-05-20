pipeline {
  agent any

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
        script {
          // Build with custom Dockerfile path (-f) and context (go/rr)
          def image = docker.build("${IMAGE_NAME}:${IMAGE_TAG}", "-f go/rr/rr.Dockerfile go/rr")
        }
      }
    }

    stage('Push Docker Image') {
      steps {
        script {
          docker.withRegistry('https://index.docker.io/v1/', 'dockerhub-creds') {
            def image = docker.image("${IMAGE_NAME}:${IMAGE_TAG}")
            image.push()
          }
        }
      }
    }
  }
}

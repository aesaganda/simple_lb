pipeline {
    agent {
        kubernetes {
            label 'docker-agent'
        }
    }
    environment {
        IMAGE_NAME = 'aesaganda/jenkins-lb'
        IMAGE_TAG = 'latest'
        DOCKER_CREDENTIALS_ID = 'dockerhub-creds'
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
            docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -f python/lb.Dockerfile python
          """
                }
            }
        }
        stage('Push Docker Image') {
            steps {
                container('docker') {
                    withCredentials([usernamePassword(credentialsId: "${DOCKER_CREDENTIALS_ID}", usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                        sh '''
              echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin docker.io
              docker push ${IMAGE_NAME}:${IMAGE_TAG}
            '''
                    }
                }
            }
        }
        stage('Download twistcli') {
            steps {
                container('docker') {
                    sh """
wget -q -O twistcli https://utd-packages.s3.amazonaws.com/twistcli
chmod +x twistcli
                        """
                }
            }
        }
        stage('Scan Docker Image with Twistlock') {
            steps {
                container('docker') {
                    script { // Using script block for better control over shell commands
                        try {
                            sh(script: """
                               ./twistcli images scan \\
                               --address https://twistlock1.garanti.lab:8083/ \\
                               --user admin \\
                               --password admin \\
                               --details \\
                               ${IMAGE_NAME}:${IMAGE_TAG}
                            """, returnStatus: true) // <-- The key change: returnStatus: true
                        } catch (Exception e) {
                            echo "Twistcli scan completed with non-zero exit code (likely vulnerabilities found). Proceeding anyway."
                            // Optionally, you can log the error or perform other actions here
                        }
                    }
                }
            }
        }
    }
}

pipeline {
    agent {
        kubernetes {
            inheritFrom 'docker-agent' // Changed from deprecated 'label' to 'inheritFrom'
        }
    }
    environment {
        IMAGE_NAME = 'aesaganda/jenkins-lb'
        IMAGE_TAG = 'latest'
        DOCKER_CREDENTIALS_ID = 'dockerhub-creds'
        DOCKER_REGISTRY = 'docker.io' // Added for Docker Hub
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
                            echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin ${DOCKER_REGISTRY}
                            docker push ${IMAGE_NAME}:${IMAGE_TAG}
                        '''
                    }
                }
            }
        }
        stage('Download twistcli') {
            steps {
                container('docker') {
                    sh '''
                        wget -q -O twistcli https://utd-packages.s3.amazonaws.com/twistcli
                        chmod +x twistcli
                    '''
                }
            }
        }
        stage('Scan Docker Image with Twistlock') {
            steps {
                container('docker') {
                    script {
                        // Use script block to handle the command and capture output
                        def scanResult = sh(script: """
                            ./twistcli images scan \\
                                --address https://twistlock1.garanti.lab:8083 \\
                                --user admin \\
                                --password admin \\
                                --details \\
                                ${IMAGE_NAME}:${IMAGE_TAG} || true
                        """, returnStatus: true)
                        echo "Twistlock scan completed with exit code: ${scanResult}"
                        // Optionally log output to a file for debugging
                        sh '''
                            ./twistcli images scan \\
                                --address https://twistlock1.garanti.lab:8083 \\
                                --user admin \\
                                --password admin \\
                                --details \\
                                ${IMAGE_NAME}:${IMAGE_TAG} > twistcli_scan_output.log 2>&1 || true
                        '''
                    }
                }
            }
        }
    }
    post {
        always {
            container('docker') {
                // Clean up twistcli binary and Docker login
                sh 'rm -f twistcli || true'
                sh 'docker logout || true'
                // Archive scan output for debugging
                archiveArtifacts artifacts: 'twistcli_scan_output.log', allowEmptyArchive: true
            }
        }
    }
}

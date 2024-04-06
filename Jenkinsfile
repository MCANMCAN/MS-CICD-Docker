pipeline {
    agent any

    environment {
        // environment variables
        DOCKER_REGISTRY_CREDENTIALS = credentials('docker-hub-credentials')
        DOCKER_IMAGE_NAME = '17059593/go-auth'
        DOCKERFILE_PATH = './Dockerfile' // Path to your Dockerfile
    }

    stages {
        stage('Checkout') {
            steps {
                // Checkout code from Git repository
                git 'https://github.com/MCANMCAN/MS-CICD-Docker.git'
            }
        }

        stage('Dockerizing Application') {
            steps {
                // Build Docker image
                script {
                    docker.build(env.DOCKER_IMAGE_NAME, '-f ' + env.DOCKERFILE_PATH + ' .')
                }
            }
        }

        // stage('Pushing Docker Image') {
        //     steps {
        //         // Push Docker image to Docker Hub
        //         script {
        //             docker.withRegistry('https://index.docker.io/v1/', env.DOCKER_REGISTRY_CREDENTIALS) {
        //                 docker.image(env.DOCKER_IMAGE_NAME).push('latest')
        //             }
        //         }
        //     }
        // }
    }
}
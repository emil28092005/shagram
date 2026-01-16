pipeline{
    agent { label 'docker' }
    options {
        skipDefaultCheckout()
    }
    environment {
        PROJECT_DIR = "/opt/shagram/shagram"
        IMAGE_NAME = "shagram"
    }
    stages{
        stage('Checkout') {
            steps {
                dir(env.PROJECT_DIR) {
                    deleteDir()
                    git url: 'https://github.com/emil28092005/shagram.git', branch: 'main'
                }
            }
        }
        stage('Build') {
            steps {
                dir(env.PROJECT_DIR) { 
                    echo 'Building shagram...'
                    sh 'docker build -t ${IMAGE_NAME}:${GIT_COMMIT} .'

                }
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'docker run --rm ${IMAGE_NAME}:${GIT_COMMIT} sqlite3 --version'
            }
        }
        stage('Deploy') {
            steps {
                dir("${PROJECT_DIR}/deploy/shagram") {
                    echo 'Deploying...'
                    sh 'printf "APP_IMAGE=%s\\n" "${IMAGE_NAME}:${GIT_COMMIT}" > .env'
                    sh 'docker compose up -d'
                    sh 'docker compose ps'
                }
            }
        }
    }
}
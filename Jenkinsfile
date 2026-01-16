pipeline{
    agent any
    stages{
        stage('Checkout') {
            steps {
                deleteDir()
                git url: 'https://github.com/emil28092005/shagram.git', branch: 'main'
            }
        }
        stage('Build') {
            steps {
                echo 'Building shagram...'
                sh 'docker build -t shagram-shagram:latest .'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing...'
                sh 'docker run --rm shagram-shagram:latest sqlite3 /app/shagram.db "SELECT 1;" || true'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying...'
                sh '''
                    docker compose down || true
                    docker compose up -d
                    docker compose ps
                '''
            }
        }
    }
}
pipeline{
    agent any
    stages{
        stage('Build') {
            agent {
                docker {
                    image 'docker:27-dind'
                    args '--privileged -v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
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
                    docker-compose down || true
                    docker-compose up -d
                    docker-compose ps
                '''
            }
        }
    }
}
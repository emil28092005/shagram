pipeline {
  agent { label 'docker-agent' }

  options {
    skipDefaultCheckout(true)
    timestamps()
  }

  environment {
    IMAGE_NAME   = "shagram"
    DEPLOY_DIR   = "/opt/shagram/shagram/deploy/shagram"
    COMPOSE_FILE = "/opt/shagram/shagram/deploy/shagram/compose.yaml"
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
        sh '''
          set -eux
          git reset --hard
          git clean -xffd
        '''
      }
    }

    stage('Build') {
      steps {
        script {
          env.GIT_SHA = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
          env.APP_IMAGE = "${IMAGE_NAME}:${env.GIT_SHA}"
        }
        sh '''
          set -eux
          docker build -t "$APP_IMAGE" .
        '''
      }
    }

    stage('Test') {
      steps {
        sh 'set -eux; go test ./...'
      }
    }

    stage('Deploy') {
      steps {
        sh '''
          set -eux
          mkdir -p "$DEPLOY_DIR"
          cat > "$DEPLOY_DIR/.env" <<EOF
APP_IMAGE=$APP_IMAGE
EOF
          docker compose -f "$COMPOSE_FILE" --project-directory "$DEPLOY_DIR" up -d --remove-orphans
        '''
      }
    }
  }
}

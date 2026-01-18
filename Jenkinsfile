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
          git status --porcelain
        '''
      }
    }

    stage('Build') {
      steps {
        script {
          def sha = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
          env.IMAGE_TAG = sha
          env.APP_IMAGE = "${IMAGE_NAME}:${sha}"
        }

        sh '''
          set -eux
          docker version
          docker build -t "$APP_IMAGE" .
        '''
      }
    }

    stage('Test') {
      steps {
        sh '''
          set -eux
          go test ./...
        '''
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
          docker compose -f "$COMPOSE_FILE" --project-directory "$DEPLOY_DIR" ps
        '''
      }
    }
  }

  post {
    always {
      echo "Done"
    }
  }
}

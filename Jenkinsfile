pipeline {
  agent { label 'docker' }
  options { skipDefaultCheckout() }

  environment {
    IMAGE_NAME   = "shagram"
    DEPLOY_DIR   = "/opt/shagram/shagram/deploy/shagram"
    COMPOSE_FILE = "/opt/shagram/shagram/deploy/shagram/compose.yaml"
  }

  stages {
    stage('Checkout') {
      steps {
        script {
          deleteDir()

          def scmVars = git url: 'https://github.com/emil28092005/shagram.git', branch: 'cicd-fix'
          env.IMAGE_TAG = scmVars.GIT_COMMIT.take(7)
        }
      }
    }

    stage('Build') {
      steps {
        sh 'docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .'
      }
    }

    stage('Test') {
      steps {
        sh 'docker run --rm ${IMAGE_NAME}:${IMAGE_TAG} sqlite3 --version'
      }
    }

    stage('Deploy') {
      steps {
        sh '''
          set -e
          printf "APP_IMAGE=%s\\n" "${IMAGE_NAME}:${IMAGE_TAG}" > "${DEPLOY_DIR}/.env"

          docker compose \
            --project-directory "${DEPLOY_DIR}" \
            -f "${COMPOSE_FILE}" \
            up -d --build

          docker compose \
            --project-directory "${DEPLOY_DIR}" \
            -f "${COMPOSE_FILE}" \
            ps
        '''
      }
    }
  }
}

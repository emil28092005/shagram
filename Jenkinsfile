pipeline {
  agent { label 'docker-agent' }

  options {
    skipDefaultCheckout(true)
    timestamps()
  }

    environment {
    IMAGE_NAME = "shagram"
    DEPLOY_DIR = "/opt/shagram/shagram/deploy/shagram"
    COMPOSE_FILE = "${WORKSPACE}/deploy/shagram/compose.yaml"
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
        sh '''
          echo "Testing..."
        '''
      }
    }
    stage('Debug env') {
        steps {
            sh '''
            echo "BRANCH_NAME=$BRANCH_NAME"
            echo "GIT_BRANCH=$GIT_BRANCH"
            git rev-parse --abbrev-ref HEAD || true
            git rev-parse HEAD
            '''
        }
    }


    stage('Deploy') {
    when {
        expression { env.GIT_BRANCH == 'origin/main' || env.GIT_BRANCH == 'main' }
    }
    steps {
        sh 'echo "deploying"'
        sh 'touch /opt/shagram/TEST_WEBHOOK.txt'
    }
    }


  }
}

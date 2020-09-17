pipeline {
  agent any
  stages {
    stage('Checkout') {
      steps {
        git(url: 'https://github.com/APoniatowski/GoSSH', branch: 'v2')
      }
    }

    stage('Go Test') {
      steps {
        echo 'Fetching dependencies and last minute tidy'
        sh 'go mod tidy'
        echo 'Running go test'
        sh 'go test -v main.go'
      }
    }

    stage('Build') {
      parallel {
        stage('Linux') {
          steps {
            echo 'Compiling GoSSH - Linux'
            sh 'GOOS=linux go build -v'
          }
        }

        stage('Windows') {
          steps {
            echo 'Compiling GoSSH - Windows'
            sh 'GOOS=windows go build -v'
          }
        }

        stage('FreeBSD') {
          steps {
            echo 'Compiling GoSSH - FreeBSD'
            sh 'GOOS=freebsd go build -o "GOSSH-bsd" -v'
          }
        }

        stage('MacOS') {
          steps {
            echo 'Compiling GoSSH - MAC'
            sh 'GOOS=darwin go build -o "GoSSH-mac" -v'
          }
        }

      }
    }

    stage('SooS') {
      steps {
        echo 'Placeholder for further post-compilation testing'
        echo 'and whatever else I can think of. I am setting up'
        echo 'a small virtual infrastructure for testing'
      }
    }

  }
  environment {
    XDG_CACHE_HOME = '/tmp'
    CGO_ENABLED = '0'
  }
}
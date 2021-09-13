pipeline {
    agent {
        label "linux"
    }
    stages {
        stage('build') {
            steps {
                sh """
                    docker build -t clean_go .
                """
            }
        }
        stage('run') {
            steps {
                sh """
                    docker run --rm clean_go
                """
            }
        }
    }
}
pipeline {
    agent {
        docker
    }
    stages {
        stage('build') {
            steps {
                echo "building stage"
                sh """
                    docker build -t clean_go .
                """
            }
        }
        stage('run') {
            steps {
                echo "running stage"
                sh """
                    docker run --rm clean_go
                """
            }
        }
    }
}
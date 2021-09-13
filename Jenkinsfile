pipeline {
    agent { 
        dockerfile true 
    }
    stages {
        stage('build') {
            steps {
                echo "building stage"
                sh """
                    sudo docker build -t clean_go .
                """
            }
        }
        stage('run') {
            steps {
                echo "running stage"
                sh """
                    sudo docker run --rm clean_go
                """
            }
        }
    }
}
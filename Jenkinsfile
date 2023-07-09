node {
   def commit_id
   stage('Preparation') {
     checkout scm
     sh "git rev-parse --short HEAD > .git/commit-id"
     commit_id = readFile('.git/commit-id').trim()
   }
   stage('test') {
     def golangTestContainer = docker.image('golang:1.20')
     golangTestContainer.pull()
     golangTestContainer.inside {
//        sh 'go test ./...'
     }
   }
   stage('test with a DB') {
     def mysql = docker.image('mysql').run("-e MYSQL_ALLOW_EMPTY_PASSWORD=yes")
     def myTestContainer = docker.image('golang:1.20')
     myTestContainer.pull()
     myTestContainer.inside("--link ${mysql.id}:mysql") {
//        sh 'go test ./...'
     }
     mysql.stop()
   }
   stage('docker build/push') {
     docker.withRegistry('https://index.docker.io/v2/', 'dockerhub') {
        def app = docker.build("coenpark/coen-chat:${commit_id}", '.').push()
     }
   }
}
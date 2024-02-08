properties([
  buildDiscarder(
    logRotator(
      artifactDaysToKeepStr: '',
      artifactNumToKeepStr: '',
      daysToKeepStr: '',
      numToKeepStr: '10'
    )
  ),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron("H H * * *")
  ])
])
node("go-arm64") {
  stage("Checkout") {
    checkout scm
  }
  stage("Init") {
    sh 'make clean init test'
  }
  stage("linux_amd64") {
    sh 'make -f Makefile.gen linux_amd64'
  }
}
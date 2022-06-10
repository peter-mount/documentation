// Build properties
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
    cron('H H * * *'),
    githubPush()
  ])
])

// Image name being created
TAG="documentation:latest"

node( 'documentation' ) {

    stage( 'Prepare' ) {
        checkout([$class: 'GitSCM', branches: [[name: '*/master']], extensions: [], userRemoteConfigs: [[credentialsId: 'peter-ssh', url: 'https://github.com/peter-mount/documentation']]])

        // Form cmd, must be done in a stage as we need the build agent's userId
        userId = sh returnStdout: true, script: "id -u"
        cmd="docker run -i --rm -u ${userId.trim()} -v ${env.WORKSPACE}:/work " + TAG + " doctool "
        sh "echo '${cmd}'"
    }

    stage( "vendors" ) {
        sh "mkdir -p themes/area51/assets/vendor"
        dir('themes/area51/assets/vendor') {
            sh "git clone https://github.com/twbs/bootstrap"
            sh "git clone https://github.com/FortAwesome/Font-Awesome"
        }
    }

    stage( "build" ) {
        sh "docker-build -t ${TAG} ."
    }

    stage( "generate") {
        sh "${cmd} -p"
    }

    stage( "pdf") {
        sh "${cmd}"
    }

    stage( "publish" ) {
        sh "rsync-web1 public/ /var/www/html"
    }
}

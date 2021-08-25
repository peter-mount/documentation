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

version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

node( 'documentation' ) {
    stage( 'Prepare' ) {
        checkout scm
    }

    stage( "Compile" ) {
        sh "docker run -i --rm " +
            "-v \$(pwd):/work " +
            "-e USERID=\$(id -u) " +
            "docker.europa.area51.dev/area51/documentation:latest " +
            "compile.sh"
    }
}

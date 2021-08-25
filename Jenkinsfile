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

    stage( 'test' ) {
        sh "pwd"
        //sh "df -h"
        sh "ls -l"
    }

}

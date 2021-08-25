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
    cron('H H * * *')
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
        //sh "df -h"
    }
}

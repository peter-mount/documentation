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

tag = "docker.europa.area51.dev/area51/documentation:latest"

cmd = "docker run -i --rm -v \$(pwd):/work -e USERID=\$(id -u) " + tag + " compile.sh "

node( 'documentation' ) {
    stage( 'Prepare' ) {
        checkout scm
        sh "rm -rf node_modules public"
        sh "docker pull " + tag
    }

    stage( "NPM" ) {
        sh cmd + "npm install"
    }

    stage( "6502" ) {
        sh cmd + "doctool -6502 content/6502/"
    }

    stage( "BBC" ) {
        sh cmd + "doctool -bbc content/bbc/"
    }

    stage( "hugo" ) {
        sh "rm -rf public"
        sh cmd + "hugo"
    }

    stage( "upload" ) {
        //sh "cd public;find . -type f -exec curl -s --user raw-upload:secret --ftp-create-dirs -T {} https://nexus.europa.area51.dev/repository/raw-hugo/" + version + "/{} \\;"
        sh "cd public;find . -type f -exec curl -s --user raw-upload:secret --upload-file \"{}\" \"https://nexus.europa.area51.dev/repository/raw-hugo/" + version + "/{}\" \\;"
    }
}

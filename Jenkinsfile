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

tag = "docker.europa.area51.dev/area51/documentation:" + version

cmd = "docker run -i --rm -v \$(pwd):/work -e USERID=\$(id -u) " + tag + " compile.sh "

node( 'documentation' ) {
    stage( 'Prepare' ) {
        checkout scm
    }

    stage( "Toolset" ) {
        sh "docker build -t " + tag + " " +
            "--build-arg prefix=docker.europa.area51.dev/library/ " +
            "--build-arg aptrepo=https://nexus.europa.area51.dev/repository/apt- " +
            "--add-host nexus.europa.area51.dev:192.168.2.4 " +
            "--build-arg npmrepo=https://nexus.europa.area51.dev/repository/npm-group/ " +
            "."
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
        sh "find public -type f -exec curl --user raw-upload:secret --ftp-create-dirs -T {} https://nexus.europa.area51.dev/repository/raw-hugo/" + version + " \\;"
    }
}

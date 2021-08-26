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

// List of books to generate PDF versions
def books = [ "bbc", "6502" ]

// List of references to generate via doctool
def references = [
    "-6502 content/6502/",
    "-bbc content/bbc/"
]

version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

tag = "docker.europa.area51.dev/area51/documentation:latest"

cmd = "docker run -i --rm -v \$(pwd):/work -e USERID=\$(id -u) " + tag + " compile.sh "

node( 'documentation' ) {
    stage( 'Prepare' ) {
        checkout scm
        sh "rm -rf pdfs node_modules public"
        sh "docker pull " + tag
    }

    stage( "npm" ) {
        sh cmd + "npm install"
    }

    stage( "reference generation" ) {
        for( reference in references ) {
            sh cmd + "doctool " + reference
        }
    }

    stage( "hugo" ) {
        sh "rm -rf public"
        sh cmd + "hugo"
    }

    stage( "PDF generation" ) {
        for( book in books ) {
            sh cmd + "generate-pdf.sh " + book
        }
    }

    stage( "upload" ) {
        //sh "cd public;find . -type f -exec curl -s --user raw-upload:secret --ftp-create-dirs -T {} https://nexus.europa.area51.dev/repository/raw-hugo/" + version + "/{} \\;"
        sh "cd public;find . -type f -exec curl -s --user raw-upload:secret --upload-file \"{}\" \"https://nexus.europa.area51.dev/repository/raw-hugo/" + version + "/{}\" \\;"
    }
}

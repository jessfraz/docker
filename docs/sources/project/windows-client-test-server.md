page_title: Set up a Windows server for client testing
page_description: How to set up a server to test Docker Windows client
page_keywords: development, inception, container, image Dockerfile, dependencies, Go, artifacts, windows


# Windows Client test server

In this section, you will learn how to setup a Windows machine for building and
testing the Docker Windows client.

Things you will need:

- Windows Server/Desktop Machine
- Java
- MSysGit
- TDM-GCC
- MinGW (tar and xz)
- Go with linux/amd64 cross-compiling
- Docker.exe

## Installing Java

1. Open the desktop, and then tap or click the Internet Explorer icon on the
taskbar.

2. Go to [java.com](http://java.com).

3. Tap or click the Free Java Download button, and then tap or click Agree and
Start Free Download. ...

4. On the notification bar, tap or click Run.


## Installing MSysGit

1. In Internet Explorer go to [msysgit.github.io](https://msysgit.github.io/).

2. Tap or click the button to download MSysGit.

## Installing GCC

1. In Internet Explorer go to
   [tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download).

2. Tap or click the lastest version of the package, install and run.

## Installing MinGW (tar and xz)

1. In Internet Explorer download the latest MinGW from
   [SourceForge](http://sourceforge.net/projects/mingw/).

2. This will install the MinGW Installation Manager.

3. In the MinGW Installation Manager, make sure you select `tar` and `xz` as
   shown below, then select Installation > Apply Changes, to install the
   selected packages.
    
    ![windows-mingw](/project/images/windows-mingw.png)

## Installing Golang with linux/amd64 cross-compiling

1. Lets make sure gcc is in your `Path`.

2. Go to Control Panel > System and Security > System, then click the Advanced
   System Settings link in the sidebar.

3. In the Advanced tab, click Enviornment Variables. Now scroll to the Path
   variable under System variables, and double-click it to see what it
   includes.

    ![windows-env-vars](/project/images/windows-env-vars.png)

4. Make sure it has at least `C:\TDM-GCC64` so we can cross compile Go.

5. In Internet Explorer, download and install the latest `.msi` installer from
   [golang.org/dl](http://golang.org/dl/).

6. After that is done installing open gitbash.

7. Enter the directory where go was installed. For basic installations, run `cd
   /c/Go/src`.

8. Run `cmd.exe`. You should now be in the `cmd.exe` prompt.

9. Run 
        
        set GOOS=linux
        set GOARCH=amd64
        make.bat

10. This will allow for cross compiling for linux/amd64. If this fails with
    cannot find `gcc` go back to step 1.

11. Add the `GOPATH` enviornment variables. Create a folder
    `C:\gopath\src\github.com\docker\docker` and clone docker into it. Set the
    `GOPATH` variable in the screen from step 3 as
    `C:\gopath;C:\gopath\src\github.com\docker\docker\vendor`.

## Install docker.exe

1. Create a directory `C:\bin`.

2. Download the latest master build of the windows client from
   [master.dockerproject.com](https://master.dockerproject.com).

3. Place this file in the directory `C:\bin` and make sure that directory is in
   your path. Go to Control Panel > System and Security > System, then click the Advanced
   System Settings link in the sidebar.

4. In the Advanced tab, click Enviornment Variables. Now scroll to the Path
   variable under System variables, and double-click it to see what it
   includes. Add `C:\bin` to the end of the variable if it is not already
   there.

    ![windows-env-vars](/project/images/windows-env-vars.png)

5. Set a `DOCKER_HOST` variable that will connect to your remote linux instance
   with docker installed. You will also want to set the values of
   `DOCKER_CLIENTONLY` to `1` and `DOCKER_TEST_HOST` to save value as
   `DOCKER_HOST`.

6. To test the installation and remote host open GitBash and run a simple
   `docker info`, or any docker command.


## Building and testing the Windows Docker Client
You should now have a windows server and remote host all set up for building and testing docker. 

Let's try it.

1. Open GitBash.

2. Enter the docker source directory we set up `cd
   /c/gopath/src/github.com/docker/docker`.

3. Run `project/make.sh binary test-integration-cli` to make the binary and run
   the integration tests.

If you make any changes just run these commands again.

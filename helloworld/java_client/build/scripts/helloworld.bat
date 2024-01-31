@rem
@rem Copyright 2015 the original author or authors.
@rem
@rem Licensed under the Apache License, Version 2.0 (the "License");
@rem you may not use this file except in compliance with the License.
@rem You may obtain a copy of the License at
@rem
@rem      https://www.apache.org/licenses/LICENSE-2.0
@rem
@rem Unless required by applicable law or agreed to in writing, software
@rem distributed under the License is distributed on an "AS IS" BASIS,
@rem WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@rem See the License for the specific language governing permissions and
@rem limitations under the License.
@rem

@if "%DEBUG%" == "" @echo off
@rem ##########################################################################
@rem
@rem  helloworld startup script for Windows
@rem
@rem ##########################################################################

@rem Set local scope for the variables with windows NT shell
if "%OS%"=="Windows_NT" setlocal

set DIRNAME=%~dp0
if "%DIRNAME%" == "" set DIRNAME=.
set APP_BASE_NAME=%~n0
set APP_HOME=%DIRNAME%..

@rem Resolve any "." and ".." in APP_HOME to make it shorter.
for %%i in ("%APP_HOME%") do set APP_HOME=%%~fi

@rem Add default JVM options here. You can also use JAVA_OPTS and HELLOWORLD_OPTS to pass JVM options to this script.
set DEFAULT_JVM_OPTS=

@rem Find java.exe
if defined JAVA_HOME goto findJavaFromJavaHome

set JAVA_EXE=java.exe
%JAVA_EXE% -version >NUL 2>&1
if "%ERRORLEVEL%" == "0" goto execute

echo.
echo ERROR: JAVA_HOME is not set and no 'java' command could be found in your PATH.
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:findJavaFromJavaHome
set JAVA_HOME=%JAVA_HOME:"=%
set JAVA_EXE=%JAVA_HOME%/bin/java.exe

if exist "%JAVA_EXE%" goto execute

echo.
echo ERROR: JAVA_HOME is set to an invalid directory: %JAVA_HOME%
echo.
echo Please set the JAVA_HOME variable in your environment to match the
echo location of your Java installation.

goto fail

:execute
@rem Setup the command line

set CLASSPATH=%APP_HOME%\lib\helloworld.jar;%APP_HOME%\lib\grpc-all-1.52.0.jar;%APP_HOME%\lib\gax-grpc-2.7.1.jar;%APP_HOME%\lib\grpc-xds-1.52.0.jar;%APP_HOME%\lib\grpc-alts-1.52.0.jar;%APP_HOME%\lib\grpc-grpclb-1.52.0.jar;%APP_HOME%\lib\grpc-services-1.52.0.jar;%APP_HOME%\lib\protobuf-java-util-4.0.0-rc-2.jar;%APP_HOME%\lib\grpc-rls-1.52.0.jar;%APP_HOME%\lib\opencensus-proto-0.2.0.jar;%APP_HOME%\lib\grpc-protobuf-1.52.0.jar;%APP_HOME%\lib\proto-google-common-protos-2.9.0.jar;%APP_HOME%\lib\protobuf-java-4.0.0-rc-2.jar;%APP_HOME%\lib\gax-2.7.1.jar;%APP_HOME%\lib\api-common-2.0.5.jar;%APP_HOME%\lib\javax.annotation-api-1.3.2.jar;%APP_HOME%\lib\netty-tcnative-2.0.45.Final-linux-x86_64.jar;%APP_HOME%\lib\grpc-auth-1.52.0.jar;%APP_HOME%\lib\grpc-netty-1.52.0.jar;%APP_HOME%\lib\grpc-okhttp-1.52.0.jar;%APP_HOME%\lib\grpc-testing-1.52.0.jar;%APP_HOME%\lib\grpc-netty-shaded-1.52.0.jar;%APP_HOME%\lib\grpc-core-1.52.0.jar;%APP_HOME%\lib\grpc-stub-1.52.0.jar;%APP_HOME%\lib\grpc-protobuf-lite-1.52.0.jar;%APP_HOME%\lib\grpc-api-1.52.0.jar;%APP_HOME%\lib\google-auth-library-oauth2-http-1.4.0.jar;%APP_HOME%\lib\google-http-client-gson-1.41.0.jar;%APP_HOME%\lib\google-http-client-1.41.0.jar;%APP_HOME%\lib\opencensus-contrib-http-util-0.28.0.jar;%APP_HOME%\lib\guava-31.1-android.jar;%APP_HOME%\lib\error_prone_annotations-2.14.0.jar;%APP_HOME%\lib\gson-2.9.0.jar;%APP_HOME%\lib\opencensus-api-0.28.0.jar;%APP_HOME%\lib\grpc-context-1.52.0.jar;%APP_HOME%\lib\jsr305-3.0.2.jar;%APP_HOME%\lib\google-auth-library-credentials-1.4.0.jar;%APP_HOME%\lib\threetenbp-1.5.0.jar;%APP_HOME%\lib\annotations-4.1.1.4.jar;%APP_HOME%\lib\animal-sniffer-annotations-1.21.jar;%APP_HOME%\lib\perfmark-api-0.25.0.jar;%APP_HOME%\lib\netty-codec-http2-4.1.79.Final.jar;%APP_HOME%\lib\netty-handler-proxy-4.1.79.Final.jar;%APP_HOME%\lib\netty-codec-http-4.1.79.Final.jar;%APP_HOME%\lib\netty-handler-4.1.79.Final.jar;%APP_HOME%\lib\netty-transport-native-unix-common-4.1.79.Final.jar;%APP_HOME%\lib\okio-1.17.5.jar;%APP_HOME%\lib\auto-value-annotations-1.9.jar;%APP_HOME%\lib\junit-4.12.jar;%APP_HOME%\lib\re2j-1.6.jar;%APP_HOME%\lib\failureaccess-1.0.1.jar;%APP_HOME%\lib\listenablefuture-9999.0-empty-to-avoid-conflict-with-guava.jar;%APP_HOME%\lib\checker-qual-3.12.0.jar;%APP_HOME%\lib\j2objc-annotations-1.3.jar;%APP_HOME%\lib\netty-codec-socks-4.1.79.Final.jar;%APP_HOME%\lib\netty-codec-4.1.79.Final.jar;%APP_HOME%\lib\netty-transport-4.1.79.Final.jar;%APP_HOME%\lib\netty-buffer-4.1.79.Final.jar;%APP_HOME%\lib\netty-resolver-4.1.79.Final.jar;%APP_HOME%\lib\netty-common-4.1.79.Final.jar;%APP_HOME%\lib\hamcrest-core-1.3.jar;%APP_HOME%\lib\conscrypt-openjdk-uber-2.5.2.jar;%APP_HOME%\lib\httpclient-4.5.13.jar;%APP_HOME%\lib\httpcore-4.4.15.jar;%APP_HOME%\lib\commons-logging-1.2.jar;%APP_HOME%\lib\commons-codec-1.11.jar


@rem Execute helloworld
"%JAVA_EXE%" %DEFAULT_JVM_OPTS% %JAVA_OPTS% %HELLOWORLD_OPTS%  -classpath "%CLASSPATH%" com.k8s_demo.HelloWorldClient %*

:end
@rem End local scope for the variables with windows NT shell
if "%ERRORLEVEL%"=="0" goto mainEnd

:fail
rem Set variable HELLOWORLD_EXIT_CONSOLE if you need the _script_ return code instead of
rem the _cmd.exe /c_ return code!
if  not "" == "%HELLOWORLD_EXIT_CONSOLE%" exit 1
exit /b 1

:mainEnd
if "%OS%"=="Windows_NT" endlocal

:omega

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/spf13/cobra"
)

var (
	file string
	info map[string]string
)

const (
	ColorReset   = "\033[0m"
	ColorKeyword = "\033[1;32m"
	ColorComment = "\033[0;36m"
	ColorString  = "\033[0;34m"
	ColorCode    = "\033[0;33m"
)

var dockerfileKeywords = []string{
	"FROM", "RUN", "CMD", "LABEL", "MAINTAINER", "EXPOSE", "ENV",
	"ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR",
	"ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
}

type DockerfileEditor struct {
}

func initializeMap() {
	info = make(map[string]string)
	info["MAINTAINER"] = "\n" + ColorKeyword + "MAINTAINER " + ColorCode + "<name>" + ColorReset + "\n" + "\nThe" + ColorCode + " MAINTAINER " + ColorReset + "instruction sets the Author field of the generated images. The" + ColorCode + " LABEL " + ColorReset + "instruction is a much more flexible version of this and you should use it instead, as it enables setting any metadata you require, and can be viewed easily, for example with " + ColorCode + "docker inspect" + ColorReset + ". To set a label corresponding to the " + ColorCode + "MAINTAINER" + ColorReset + " field you could use:\n\n" + ColorKeyword + "LABEL" + ColorCode + " org.opencontainers.image.authors=\"SvenDowideit@home.org.au\"" + "\n" + ColorReset + "\nThis will then be visible from " + ColorCode + "docker inspect" + ColorReset + " with the other labels."
	info["EXPOSE"] = "\n" + ColorKeyword + "EXPOSE " + ColorCode + "<port> [<port>/<protocol>...]" + ColorReset +
		"\n\nThe " + ColorCode + "EXPOSE " + ColorReset +
		"instruction informs Docker that the container listens on the specified network ports at runtime. You can specify whether the port listens on TCP or UDP, and the default is TCP if you don't specify a protocol.\n\nThe " +
		ColorCode + "EXPOSE " + ColorReset +
		"instruction doesn't actually publish the port. It functions as a type of documentation between the person who builds the image and the person who runs the container, about which ports are intended to be published. To publish the port when running the container, use the " +
		ColorCode + "-p" + ColorReset + " flag on " +
		ColorCode + "docker run" + ColorReset +
		" to publish and map one or more ports, or the " +
		ColorCode + "-P " + ColorReset +
		"flag to publish all exposed ports and map them to high-order ports.\n\nBy default, " +
		ColorCode + "EXPOSE" + ColorReset +
		" assumes TCP. You can also specify UDP:\n\n" +
		ColorKeyword + "EXPOSE " + ColorCode + "80/udp\n\n" +
		ColorReset +
		"To expose on both TCP and UDP, include two lines:\n\n" +
		ColorKeyword + "EXPOSE " + ColorCode + "80/tcp\n" +
		ColorKeyword + "EXPOSE " + ColorCode + "80/udp\n\n" +
		ColorReset +
		"In this case, if you use " +
		ColorCode + "-P" + ColorReset + " with " +
		ColorCode + "docker run" + ColorReset +
		", the port will be exposed once for TCP and once for UDP. Remember that " +
		ColorCode + "-P " + ColorReset +
		"uses an ephemeral high-ordered host port on the host, so TCP and UDP doesn't use the same port.\n\nRegardless of the " +
		ColorKeyword + "EXPOSE " + ColorReset +
		"settings, you can override them at runtime by using the " +
		ColorCode + "-p " + ColorReset +
		"flag. For example\n\n " +
		ColorCode + "docker run -p 80:80/tcp -p 80:80/udp ...\n\n" +
		ColorReset +
		"To set up port redirection on the host system, see using the " +
		ColorCode + "-P" + ColorReset +
		" flag. The " +
		ColorCode + "docker network " + ColorReset +
		"command supports creating networks for communication among containers without the need to expose or publish specific ports, because the containers connected to the network can communicate with each other over any port.\n\n"
	info["ENV"] = "\n" + ColorKeyword + "ENV " + ColorCode + "<key>=<value> [<key>=<value>...]" + ColorReset +
		"\n\nThe " + ColorCode + "ENV " + ColorReset +
		"instruction sets the environment variable " + ColorCode + "<key>" + ColorReset +
		" to the value " + ColorCode + "<value>" + ColorReset +
		". This value will be in the environment for all subsequent instructions in the build stage and can be replaced inline in many as well. The value will be interpreted for other environment variables, so quote characters will be removed if they are not escaped. Like command line parsing, quotes and backslashes can be used to include spaces within values.\n\nExample:\n\n" +
		ColorKeyword + "ENV " + ColorCode + "MY_NAME=\"John Doe\"\n" +
		ColorKeyword + "ENV " + ColorCode + "MY_DOG=Rex\\ The\\ Dog\n" +
		ColorKeyword + "ENV " + ColorCode + "MY_CAT=fluffy\n\n" +
		ColorReset +
		"The " + ColorCode + "ENV " + ColorReset +
		"instruction allows for multiple " + ColorCode + "<key>=<value> " + ColorReset +
		"... variables to be set at one time, and the example below will yield the same net results in the final image:\n\n" +
		ColorKeyword + "ENV " + ColorCode + "MY_NAME=\"John Doe\" MY_DOG=Rex\\ The\\ Dog \\\n" +
		"    MY_CAT=fluffy\n\n" +
		ColorReset +
		"The environment variables set using " + ColorCode + "ENV " + ColorReset +
		"will persist when a container is run from the resulting image. You can view the values using " +
		ColorCode + "docker inspect" + ColorReset +
		", and change them using " +
		ColorCode + "docker run --env <key>=<value>" + ColorReset + ".\n\nA stage inherits any environment variables that were set using " +
		ColorCode + "ENV " + ColorReset +
		"by its parent stage or any ancestor. Refer to the multi-stage builds section in the manual for more information.\n\nEnvironment variable persistence can cause unexpected side effects. For example, setting " +
		ColorCode + "ENV DEBIAN_FRONTEND=noninteractive" + ColorReset +
		" changes the behavior of " +
		ColorCode + "apt-get" + ColorReset +
		", and may confuse users of your image.\n\nIf an environment variable is only needed during build, and not in the final image, consider setting a value for a single command instead:\n\n" +
		ColorCode + "RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y ...\n\n" +
		ColorReset +
		"Or using " +
		ColorCode + "ARG" + ColorReset +
		", which is not persisted in the final image:\n\n" +
		ColorKeyword + "ARG " + ColorCode + "DEBIAN_FRONTEND=noninteractive\n" +
		ColorCode + "RUN apt-get update && apt-get install -y ...\n\n" +
		ColorReset
	info["CMD"] = "\n" + ColorKeyword + "CMD " + ColorCode + "<command>" + ColorReset +
		"\n\nThe " + ColorCode + "CMD " + ColorReset +
		"instruction sets the command to be executed when running a container from an image.\n\nYou can specify " +
		ColorCode + "CMD " + ColorReset +
		"instructions using shell or exec forms:\n\n" +
		ColorKeyword + "CMD " + ColorCode + "[\"executable\",\"param1\",\"param2\"]" + ColorReset + " (exec form)\n" +
		ColorKeyword + "CMD " + ColorCode + "[\"param1\",\"param2\"]" + ColorReset + " (exec form, as default parameters to ENTRYPOINT)\n" +
		ColorKeyword + "CMD " + ColorCode + "command param1 param2" + ColorReset + " (shell form)\n\n" +
		"There can only be one " + ColorCode + "CMD " + ColorReset +
		"instruction in a Dockerfile. If you list more than one " +
		ColorCode + "CMD" + ColorReset +
		", only the last one takes effect.\n\nThe purpose of a " +
		ColorCode + "CMD " + ColorReset +
		"is to provide defaults for an executing container. These defaults can include an executable, or they can omit the executable, in which case you must specify an " +
		ColorCode + "ENTRYPOINT " + ColorReset +
		"instruction as well.\n\nIf you would like your container to run the same executable every time, then you should consider using " +
		ColorCode + "ENTRYPOINT " + ColorReset +
		"in combination with " +
		ColorCode + "CMD" + ColorReset +
		". See " +
		ColorCode + "ENTRYPOINT" + ColorReset +
		". If the user specifies arguments to " +
		ColorCode + "docker run" + ColorReset +
		" then they will override the default specified in " +
		ColorCode + "CMD" + ColorReset +
		", but still use the default " +
		ColorCode + "ENTRYPOINT" + ColorReset +
		".\n\nIf " +
		ColorCode + "CMD " + ColorReset +
		"is used to provide default arguments for the " +
		ColorCode + "ENTRYPOINT " + ColorReset +
		"instruction, both the " +
		ColorCode + "CMD " + ColorReset +
		"and " +
		ColorCode + "ENTRYPOINT " + ColorReset +
		"instructions should be specified in the exec form.\n\nNote\nDon't confuse " +
		ColorCode + "RUN " + ColorReset +
		"with " +
		ColorCode + "CMD" + ColorReset +
		". " +
		ColorCode + "RUN " + ColorReset +
		"actually runs a command and commits the result; " +
		ColorCode + "CMD" + ColorReset +
		" doesn't execute anything at build time, but specifies the intended command for the image.\n\n"
	info["LABEL"] = "\n" + ColorKeyword + "LABEL " + ColorCode + "<key>=<value> [<key>=<value>...]" + ColorReset +
		"\n\nThe " + ColorCode + "LABEL " + ColorReset +
		"instruction adds metadata to an image. A " +
		ColorCode + "LABEL " + ColorReset +
		"is a key-value pair. To include spaces within a " +
		ColorCode + "LABEL " + ColorReset +
		"value, use quotes and backslashes as you would in command-line parsing. A few usage examples:\n\n" +
		ColorKeyword + "LABEL " + ColorCode + "\"com.example.vendor\"=\"ACME Incorporated\"\n" +
		ColorKeyword + "LABEL " + ColorCode + "com.example.label-with-value=\"foo\"\n" +
		ColorKeyword + "LABEL " + ColorCode + "version=\"1.0\"\n" +
		ColorKeyword + "LABEL " + ColorCode + "description=\"This text illustrates \\\nthat label-values can span multiple lines.\"\n\n" +
		ColorReset +
		"An image can have more than one label. You can specify multiple labels on a single line. Prior to " +
		ColorCode + "Docker 1.10" + ColorReset +
		", this decreased the size of the final image, but this is no longer the case. You may still choose to specify multiple labels in a single instruction, in one of the following two ways:\n\n" +
		ColorKeyword + "LABEL " + ColorCode + "multi.label1=\"value1\" multi.label2=\"value2\" other=\"value3\"\n\n" +
		ColorKeyword + "LABEL " + ColorCode + "multi.label1=\"value1\" \\\n" +
		"      multi.label2=\"value2\" \\\n" +
		"      other=\"value3\"\n\n" +
		ColorReset +
		"Note\nBe sure to use double quotes and not single quotes. Particularly when you are using string interpolation (e.g. " +
		ColorCode + "LABEL example=\"foo-$ENV_VAR\"" + ColorReset +
		"), single quotes will take the string as is without unpacking the variable's value.\n\n" +
		"Labels included in base images (images in the " +
		ColorCode + "FROM " + ColorReset +
		"line) are inherited by your image. If a label already exists but with a different value, the most-recently-applied value overrides any previously-set value.\n\n" +
		"To view an image's labels, use the " +
		ColorCode + "docker image inspect" + ColorReset +
		" command. You can use the " +
		ColorCode + "--format" + ColorReset +
		" option to show just the labels;\n\n " +
		ColorCode + "docker image inspect --format='{{json .Config.Labels}}' <image>\n\n" +
		ColorReset
	info["ADD"] = "\n" + ColorKeyword + "ADD " + ColorCode + "[OPTIONS] <src> ... <dest>" + ColorReset +
		"\n" + ColorKeyword + "ADD " + ColorCode + "[OPTIONS] [\"<src>\", ... \"<dest>\"]" + ColorReset +
		"\n\nThe " + ColorCode + "ADD " + ColorReset +
		"instruction has two forms. The latter form is required for paths containing whitespace.\n\nThe available " +
		ColorCode + "[OPTIONS]" + ColorReset + " are:\n\n" +
		ColorCode + "--keep-git-dir" + ColorReset + " (1.1)\n" +
		ColorCode + "--checksum" + ColorReset + " (1.6)\n" +
		ColorCode + "--chown" + ColorReset + "\n" +
		ColorCode + "--chmod" + ColorReset + " (1.2)\n" +
		ColorCode + "--link" + ColorReset + " (1.4)\n" +
		ColorCode + "--exclude" + ColorReset + " (1.19)\n\n" +
		"The " + ColorCode + "ADD " + ColorReset +
		"instruction copies new files or directories from " +
		ColorCode + "<src>" + ColorReset +
		" and adds them to the filesystem of the image at the path " +
		ColorCode + "<dest>" + ColorReset +
		". Files and directories can be copied from the build context, a remote URL, or a Git repository.\n\nThe " +
		ColorCode + "ADD " + ColorReset +
		"and " +
		ColorCode + "COPY " + ColorReset +
		"instructions are functionally similar, but serve slightly different purposes.\n\nSource\nYou can specify multiple source files or directories with " +
		ColorCode + "ADD" + ColorReset +
		". The last argument must always be the destination. For example:\n\n" +
		ColorKeyword + "ADD " + ColorCode + "file1.txt file2.txt /usr/src/things/\n\n" +
		ColorReset +
		"If you specify multiple source files, either directly or using a wildcard, then the destination must be a directory (must end with a slash " +
		ColorCode + "/" + ColorReset +
		").\n\nTo add files from a remote location, you can specify a URL or the address of a Git repository as the source. For example:\n\n" +
		ColorKeyword + "ADD " + ColorCode + "https://example.com/archive.zip /usr/src/things/\n" +
		ColorKeyword + "ADD " + ColorCode + "git@github.com:user/repo.git /usr/src/things/\n\n" +
		ColorReset +
		"BuildKit detects the type of " +
		ColorCode + "<src>" + ColorReset +
		" and processes it accordingly.\n\n" +
		"If " + ColorCode + "<src>" + ColorReset + " is a local file or directory, the contents are copied.\n" +
		"If " + ColorCode + "<src>" + ColorReset + " is a local tar archive, it is decompressed and extracted.\n" +
		"If " + ColorCode + "<src>" + ColorReset + " is a URL, the contents are downloaded.\n" +
		"If " + ColorCode + "<src>" + ColorReset + " is a Git repository, the repository is cloned.\n\n" +
		"Adding files from the build context\nAny relative or local path that doesn't begin with " +
		ColorCode + "http://, https://, or git@" + ColorReset +
		" is considered a local file path.\n\nTrailing slashes in the source path are disregarded. Directory contents are copied, not the directory itself. Conflicts are resolved in favor of the content being added, except when copying a directory onto an existing file.\n\n" +
		"If the source is a file, file permissions and metadata are preserved. Errors occur if file/directory name conflicts exist.\n\n" +
		"If you pass a Dockerfile through stdin to the build, there is no build context. In this case, you can only use " +
		ColorCode + "ADD " + ColorReset +
		"to copy remote files.\n\nPattern matching\nFor local files, each " +
		ColorCode + "<src>" + ColorReset +
		" may contain wildcards using Go's filepath.Match rules.\n\n" +
		ColorKeyword + "ADD " + ColorCode + "*.png /dest/\n" +
		ColorKeyword + "ADD " + ColorCode + "index.?s /dest/\n\n" +
		ColorReset +
		"When adding paths with special characters, escape them using Golang rules:\n\n" +
		ColorKeyword + "ADD " + ColorCode + "arr[[]0].txt /dest/\n\n" +
		ColorReset +
		"Adding local tar archives\nLocal tar archives in recognized formats are unpacked by default. Extraction behavior matches " +
		ColorCode + "tar -x" + ColorReset +
		". Conflicts are resolved in favor of added content.\n\nNote\nFile type detection is based on file contents, not filename.\n\n" +
		"Adding files from a URL\nRemote files receive permissions " +
		ColorCode + "600" + ColorReset +
		". If the destination ends with a slash, the filename is inferred from the URL.\n\n" +
		"If authentication is required, use " +
		ColorCode + "RUN wget" + ColorReset +
		" or " +
		ColorCode + "RUN curl" + ColorReset +
		".\n\n" +
		"Adding files from a Git repository\nGit repositories are cloned into the destination. URL fragments may specify branches, tags, commits, or subdirectories.\n\n" +
		ColorKeyword + "ADD " + ColorCode + "git@github.com:moby/buildkit.git#v0.14.1:docs /buildkit-docs\n\n" +
		ColorReset +
		"Files default to permissions " +
		ColorCode + "644" + ColorReset +
		", executables " +
		ColorCode + "755" + ColorReset +
		", directories " +
		ColorCode + "755" + ColorReset +
		". SSH access requires passing credentials using the " +
		ColorCode + "--ssh" + ColorReset +
		" flag.\n\nDestination\nAbsolute paths begin with " +
		ColorCode + "/" + ColorReset +
		". Trailing slashes are significant.\n\n" +
		ColorKeyword + "ADD " + ColorCode + "test.txt /abs/\n\n" +
		ColorReset +
		"Relative destinations are relative to " +
		ColorCode + "WORKDIR" + ColorReset +
		". Missing directories are created automatically.\n\n" +
		ColorKeyword + "ADD " + ColorCode + "--keep-git-dir=true https://github.com/moby/buildkit.git#v0.10.1 /buildkit\n\n" +
		ColorReset +
		"The " + ColorCode + "--checksum" + ColorReset +
		" flag verifies remote HTTP(S) resources using " +
		ColorCode + "sha256" + ColorReset +
		".\n\n" +
		ColorKeyword + "ADD " + ColorCode + "--checksum=sha256:<hash> https://example.com/file.tar.gz /\n\n" +
		ColorReset +
		"The " + ColorCode + "--unpack" + ColorReset +
		" flag controls automatic extraction of tar archives.\n\n" +
		ColorKeyword + "ADD " + ColorCode + "--unpack=true https://example.com/archive.tar.gz /download\n" +
		ColorKeyword + "ADD " + ColorCode + "--unpack=false my-archive.tar.gz .\n\n" +
		ColorReset
	info["COPY"] = "\n" + ColorKeyword + "COPY " + ColorCode + "[OPTIONS] <src> ... <dest>" + ColorReset +
		"\n" + ColorKeyword + "COPY " + ColorCode + "[OPTIONS] [\"<src>\", ... \"<dest>\"]" + ColorReset +
		"\n\nThe " + ColorCode + "COPY " + ColorReset +
		"instruction has two forms. The latter form is required for paths containing whitespace.\n\nThe available " +
		ColorCode + "[OPTIONS]" + ColorReset + " are:\n\n" +
		ColorCode + "--from" + ColorReset + "\n" +
		ColorCode + "--chown" + ColorReset + "\n" +
		ColorCode + "--chmod" + ColorReset + " (1.2)\n" +
		ColorCode + "--link" + ColorReset + " (1.4)\n" +
		ColorCode + "--parents" + ColorReset + " (1.20)\n" +
		ColorCode + "--exclude" + ColorReset + " (1.19)\n\n" +
		"The " + ColorCode + "COPY " + ColorReset +
		"instruction copies new files or directories from " +
		ColorCode + "<src>" + ColorReset +
		" and adds them to the filesystem of the image at the path " +
		ColorCode + "<dest>" + ColorReset +
		". Files and directories can be copied from the build context, build stage, named context, or an image.\n\nThe " +
		ColorCode + "ADD " + ColorReset +
		"and " +
		ColorCode + "COPY " + ColorReset +
		"instructions are functionally similar, but serve slightly different purposes.\n\nSource\nYou can specify multiple source files or directories with " +
		ColorCode + "COPY" + ColorReset +
		". The last argument must always be the destination. For example:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "file1.txt file2.txt /usr/src/things/\n\n" +
		ColorReset +
		"If you specify multiple source files, either directly or using a wildcard, then the destination must be a directory (must end with a slash " +
		ColorCode + "/" + ColorReset +
		").\n\n" +
		ColorCode + "COPY --from=<name> " + ColorReset +
		"lets you specify the source location to be a build stage, context, or image. Example:\n\n" +
		ColorKeyword + "FROM " + ColorCode + "golang AS build\n" +
		ColorKeyword + "WORKDIR " + ColorCode + "/app\n" +
		ColorKeyword + "RUN " + ColorCode + "--mount=type=bind,target=. go build -o /myapp ./cmd\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--from=build /myapp /usr/bin/\n\n" +
		ColorReset +
		"Copying from the build context\nPaths are interpreted as relative to the root of the build context. Parent directory navigation is removed automatically. Trailing slashes in source paths are disregarded.\n\n" +
		"If the source is a directory, the contents are copied (not the directory itself). Conflicts are resolved in favor of the added content, except when copying a directory onto an existing file.\n\n" +
		"If the source is a file, metadata and permissions are preserved. Errors occur on name conflicts with existing directories.\n\n" +
		"If there is no build context (Dockerfile via stdin), " +
		ColorCode + "COPY " + ColorReset +
		"can only copy from stages, named contexts, or images using " +
		ColorCode + "--from" + ColorReset +
		".\n\n" +
		"When using a Git repository as the build context, default permissions are " +
		ColorCode + "644" + ColorReset +
		" for files, " +
		ColorCode + "755" + ColorReset +
		" for executables and directories.\n\nPattern matching\nFor local files, each " +
		ColorCode + "<src>" + ColorReset +
		" may contain wildcards using Go's filepath.Match rules.\n\n" +
		ColorKeyword + "COPY " + ColorCode + "*.png /dest/\n" +
		ColorKeyword + "COPY " + ColorCode + "index.?s /dest/\n\n" +
		ColorReset +
		"Paths with special characters must be escaped using Golang rules:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "arr[[]0].txt /dest/\n\n" +
		ColorReset +
		"Destination\nAbsolute paths begin with " +
		ColorCode + "/" + ColorReset +
		". Trailing slashes are significant.\n\n" +
		ColorKeyword + "COPY " + ColorCode + "test.txt /abs/\n\n" +
		ColorReset +
		"Relative destinations are resolved relative to " +
		ColorCode + "WORKDIR" + ColorReset +
		". Missing directories are created automatically.\n\n" +
		ColorKeyword + "WORKDIR " + ColorCode + "/usr/src/app\n" +
		ColorKeyword + "COPY " + ColorCode + "test.txt rel/\n\n" +
		ColorReset +
		"If the source is a file and the destination doesn't end with a slash, the destination path becomes the filename.\n\n" +
		ColorCode + "COPY --from" + ColorReset +
		"\nBy default, " + ColorCode + "COPY " + ColorReset +
		"copies from the build context. Using " +
		ColorCode + "--from" + ColorReset +
		", you can copy from images, stages, or named contexts:\n\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine AS build\n" +
		ColorKeyword + "COPY " + ColorCode + ". .\n" +
		ColorKeyword + "RUN " + ColorCode + "apk add clang\n" +
		ColorKeyword + "RUN " + ColorCode + "clang -o /hello hello.c\n\n" +
		ColorKeyword + "FROM " + ColorCode + "scratch\n" +
		ColorKeyword + "COPY " + ColorCode + "--from=build /hello /\n\n" +
		ColorReset +
		"You can also copy from images:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--from=nginx:latest /etc/nginx/nginx.conf /nginx.conf\n\n" +
		ColorReset +
		"The source path for " + ColorCode + "COPY --from" + ColorReset +
		" is always resolved from filesystem root of the source image or stage.\n\n" +
		ColorCode + "COPY --chown --chmod" + ColorReset +
		"\nOnly octal notation is supported for permissions. These flags work only for Linux containers.\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--chown=55:mygroup files* /somedir/\n" +
		ColorKeyword + "COPY " + ColorCode + "--chown=bin files* /somedir/\n" +
		ColorKeyword + "COPY " + ColorCode + "--chown=1 files* /somedir/\n" +
		ColorKeyword + "COPY " + ColorCode + "--chown=10:11 files* /somedir/\n" +
		ColorKeyword + "COPY " + ColorCode + "--chown=myuser:mygroup --chmod=644 files* /somedir/\n\n" +
		ColorReset +
		"User and group name resolution relies on " +
		ColorCode + "/etc/passwd" + ColorReset +
		" and " +
		ColorCode + "/etc/group" + ColorReset +
		". Numeric IDs require no lookup.\n\n" +
		ColorCode + "COPY --chmod" + ColorReset +
		" supports variable interpolation in Dockerfile syntax 1.10+:\n\n" +
		ColorKeyword + "ARG " + ColorCode + "MODE=440\n" +
		ColorKeyword + "COPY " + ColorCode + "--chmod=$MODE . .\n\n" +
		ColorReset +
		ColorCode + "COPY --link" + ColorReset +
		"\nThis flag creates a linked layer allowing better cache reuse. Source files are copied into an empty destination directory and linked on top of previous layers.\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--link /foo /bar\n\n" +
		ColorReset +
		"Using " + ColorCode + "--link" + ColorReset +
		" enables efficient rebasing and cache reuse, especially in multi-stage builds.\n\n" +
		"Incompatibilities with " + ColorCode + "--link=false" + ColorReset +
		": COPY/ADD cannot read files from previous layers and cannot follow destination symlinks.\n\n" +
		ColorCode + "COPY --parents" + ColorReset +
		"\nPreserves parent directories of source paths:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "./x/a.txt ./y/a.txt /no_parents/\n" +
		ColorKeyword + "COPY " + ColorCode + "--parents ./x/a.txt ./y/a.txt /parents/\n\n" +
		ColorReset +
		"Selective parent preservation is possible using " +
		ColorCode + "./" + ColorReset +
		" markers in paths.\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--parents ./x/./y/*.txt /parents/\n\n" +
		ColorReset +
		"With " + ColorCode + "--parents" + ColorReset +
		", BuildKit can pack multiple COPY instructions together efficiently.\n\n" +
		ColorCode + "COPY --exclude" + ColorReset +
		"\nThe " + ColorCode + "--exclude" + ColorReset +
		" flag specifies paths to exclude, supporting wildcards:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--exclude=*.txt hom* /mydir/\n\n" +
		ColorReset +
		"Multiple excludes can be specified:\n\n" +
		ColorKeyword + "COPY " + ColorCode + "--exclude=*.txt --exclude=*.md hom* /mydir/\n\n" +
		ColorReset
	info["ENTRYPOINT"] = "\n" + ColorKeyword + "ENTRYPOINT " + ColorCode + "[\"executable\", \"param1\", \"param2\"]" + ColorReset +
		"\n\nAn " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"allows you to configure a container that will run as an executable.\n\n" +
		ColorKeyword + "ENTRYPOINT " + ColorReset +
		"has two possible forms:\n\n" +
		"The exec form, which is the preferred form:\n\n" +
		ColorKeyword + "ENTRYPOINT " + ColorCode + "[\"executable\", \"param1\", \"param2\"]\n\n" +
		ColorReset +
		"The shell form:\n\n" +
		ColorKeyword + "ENTRYPOINT " + ColorCode + "command param1 param2\n\n" +
		ColorReset +
		"For more information about the different forms, see Shell and exec form.\n\n" +
		"The following command starts a container from the nginx image with its default content, listening on port 80:\n\n " +
		ColorCode + "docker run -i -t --rm -p 80:80 nginx\n\n" +
		ColorReset +
		"Command line arguments to " + ColorCode + "docker run <image>" + ColorReset +
		" will be appended after all elements in an exec form " +
		ColorKeyword + "ENTRYPOINT " + ColorReset +
		", and will override all elements specified using " +
		ColorKeyword + "CMD" + ColorReset + ".\n\n" +
		"This allows arguments to be passed to the entry point, i.e., " +
		ColorCode + "docker run <image> -d" + ColorReset +
		" will pass the " + ColorCode + "-d" + ColorReset +
		" argument to the entry point. You can override the " +
		ColorKeyword + "ENTRYPOINT " + ColorReset +
		"instruction using the " + ColorCode + "docker run --entrypoint" + ColorReset +
		" flag.\n\n" +
		"The shell form of " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"prevents any " + ColorKeyword + "CMD " + ColorReset +
		"command line arguments from being used. It also starts your ENTRYPOINT as a subcommand of " +
		ColorCode + "/bin/sh -c" + ColorReset +
		", which does not pass signals. This means that the executable will not be the container's PID 1, and will not receive Unix signals.\n\n" +
		"In this case, your executable doesn't receive a " +
		ColorCode + "SIGTERM" + ColorReset +
		" from " + ColorCode + "docker stop <container>" + ColorReset + ".\n\n" +
		"Only the last " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"instruction in the Dockerfile will have an effect.\n\n" +
		ColorKeyword + "Exec form ENTRYPOINT example\n\n" + ColorReset +
		"You can use the exec form of " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"to set fairly stable default commands and arguments and then use either form of " +
		ColorKeyword + "CMD " + ColorReset +
		"to set additional defaults that are more likely to be changed.\n\n" +
		ColorCode +
		"FROM ubuntu\n" +
		"ENTRYPOINT [\"top\", \"-b\"]\n" +
		"CMD [\"-c\"]\n\n" +
		ColorReset +
		"When you run the container, you can see that " +
		ColorCode + "top" + ColorReset +
		" is the only process and runs as PID 1.\n\n" +
		"Understand how " + ColorKeyword + "CMD " + ColorReset +
		"and " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"interact\n\n" +
		"Both " + ColorKeyword + "CMD " + ColorReset +
		"and " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"instructions define what command gets executed when running a container. There are a few rules that describe their cooperation:\n\n" +
		"- Dockerfile should specify at least one of " + ColorKeyword + "CMD " + ColorReset +
		"or " + ColorKeyword + "ENTRYPOINT\n" + ColorReset +
		"- " + ColorKeyword + "ENTRYPOINT " + ColorReset +
		"should be defined when using the container as an executable.\n" +
		"- " + ColorKeyword + "CMD " + ColorReset +
		"should be used as a way of defining default arguments for an ENTRYPOINT command or for executing an ad-hoc command in a container.\n" +
		"- " + ColorKeyword + "CMD " + ColorReset +
		"will be overridden when running the container with alternative arguments.\n\n" +
		"Note: If " + ColorKeyword + "CMD " + ColorReset +
		"is defined from the base image, setting " +
		ColorKeyword + "ENTRYPOINT " + ColorReset +
		"will reset CMD to an empty value. In this scenario, CMD must be defined in the current image to have a value.\n\n"
	info["VOLUME"] = "\n" + ColorKeyword + "VOLUME " + ColorCode + "[\"/data\"]" + ColorReset +
		"\n\nThe " + ColorKeyword + "VOLUME " + ColorReset +
		"instruction creates a mount point with the specified name and marks it as holding externally mounted volumes from the native host or other containers. The value can be a JSON array, such as " +
		ColorCode + "VOLUME [\"/var/log/\"]" + ColorReset +
		", or a plain string with multiple arguments, such as " +
		ColorCode + "VOLUME /var/log" + ColorReset +
		" or " +
		ColorCode + "VOLUME /var/log /var/db" + ColorReset +
		". For more information, examples, and mounting instructions via the Docker client, refer to the Share Directories via Volumes documentation.\n\n" +
		"The " + ColorCode + "docker run" + ColorReset +
		" command initializes the newly created volume with any data that exists at the specified location within the base image. For example, consider the following Dockerfile snippet:\n\n" +
		ColorCode +
		"FROM ubuntu\n" +
		"RUN mkdir /myvol\n" +
		"RUN echo \"hello world\" > /myvol/greeting\n" +
		"VOLUME /myvol\n\n" +
		ColorReset +
		"This Dockerfile results in an image that causes " +
		ColorCode + "docker run" + ColorReset +
		" to create a new mount point at " +
		ColorCode + "/myvol" + ColorReset +
		" and copy the greeting file into the newly created volume.\n\n" +
		ColorKeyword + "Notes about specifying volumes\n\n" + ColorReset +
		"Keep the following things in mind about volumes in the Dockerfile.\n\n" +
		"Volumes on Windows-based containers: When using Windows-based containers, the destination of a volume inside the container must be one of:\n\n" +
		"- a non-existing or empty directory\n" +
		"- a drive other than " + ColorCode + "C:" + ColorReset + "\n\n" +
		"Changing the volume from within the Dockerfile: If any build steps change the data within the volume after it has been declared, those changes will be discarded when using the legacy builder. When using " +
		ColorCode + "BuildKit" + ColorReset +
		", the changes will instead be kept.\n\n" +
		"JSON formatting: The list is parsed as a JSON array. You must enclose words with double quotes " +
		ColorCode + "(\"\")" + ColorReset +
		" rather than single quotes " +
		ColorCode + "('')" + ColorReset +
		".\n\n" +
		"The host directory is declared at container run-time: The host directory (the mountpoint) is, by its nature, host-dependent. This is to preserve image portability, since a given host directory can't be guaranteed to be available on all hosts. For this reason, you can't mount a host directory from within the Dockerfile. The " +
		ColorKeyword + "VOLUME " + ColorReset +
		"instruction does not support specifying a host-dir parameter. You must specify the mountpoint when you create or run the container.\n\n"
	info["USER"] = "\n" + ColorKeyword + "USER " + ColorCode + "<user>[:<group>]" + ColorReset +
		"\n" + ColorKeyword + "USER " + ColorCode + "UID[:GID]" + ColorReset +
		"\n\n" +
		"The " + ColorKeyword + "USER " + ColorReset +
		"instruction sets the user name (or UID) and optionally the user group (or GID) to use as the default user and group for the remainder of the current stage. The specified user is used for " +
		ColorCode + "RUN " + ColorReset +
		"instructions and at runtime, runs the relevant " +
		ColorCode + "ENTRYPOINT " + ColorReset +
		"and " +
		ColorCode + "CMD " + ColorReset +
		"commands.\n\n" +
		"Note that when specifying a group for the user, the user will have only the specified group membership. Any other configured group memberships will be ignored.\n\n" +
		ColorKeyword + "Warning\n" + ColorReset +
		"When the user doesn't have a primary group then the image (or the next instructions) will be run with the " +
		ColorCode + "root " + ColorReset +
		"group.\n\n" +
		"On Windows, the user must be created first if it's not a built-in account. This can be done with the " +
		ColorCode + "net user " + ColorReset +
		"command called as part of a Dockerfile.\n\n" +
		ColorCode +
		"FROM microsoft/windowsservercore\n" +
		"# Create Windows user in the container\n" +
		"RUN net user /add patrick\n" +
		"# Set it for subsequent commands\n" +
		"USER patrick\n" +
		ColorReset
	info["COPY"] = "\n" + ColorKeyword + "COPY " + ColorCode + "[OPTIONS] <src> ... <dest>" + ColorReset + "\n" +
		ColorKeyword + "COPY " + ColorCode + "[OPTIONS] [\"<src>\", ... \"<dest>\"]" + ColorReset + "\n\n" +
		"The available " + ColorCode + "[OPTIONS]" + ColorReset + " are:\n\n" +
		ColorCode + "--from" + ColorReset + "\n" +
		ColorCode + "--chown" + ColorReset + "\n" +
		ColorCode + "--chmod" + ColorReset + " 1.2\n" +
		ColorCode + "--link" + ColorReset + " 1.4\n" +
		ColorCode + "--parents" + ColorReset + " 1.20\n" +
		ColorCode + "--exclude" + ColorReset + " 1.19\n\n" +
		"The " + ColorCode + "COPY " + ColorReset + "instruction copies new files or directories from " + ColorCode + "<src>" + ColorReset + " and adds them to the filesystem of the image at the path " + ColorCode + "<dest>" + ColorReset + ". Files and directories can be copied from the build context, build stage, named context, or an image.\n\n" +
		"The " + ColorCode + "ADD " + ColorReset + "and " + ColorCode + "COPY " + ColorReset + "instructions are functionally similar, but serve slightly different purposes.\n\n" +
		ColorKeyword + "Source\n" + ColorReset +
		"You can specify multiple source files or directories with " + ColorCode + "COPY" + ColorReset + ". The last argument must always be the destination. For example:\n\n" +
		ColorKeyword + "COPY file1.txt file2.txt /usr/src/things/" + ColorReset + "\n\n" +
		"Multiple source files must have a directory as the destination (ending with /).\n\n" +
		ColorKeyword + "COPY --from=build /myapp /usr/bin/" + ColorReset + "\n\n" +
		"Copying from the build context interprets paths relative to the root of the context. Leading slashes or ../ are removed, trailing slashes ignored. Directories copy only their contents, preserving metadata. Conflicts overwrite existing files unless a directory conflicts with a file.\n\n" +
		ColorKeyword + "Pattern matching\n" + ColorReset +
		"Each <src> may contain wildcards using Go's filepath.Match rules. Example:\n\n" +
		ColorKeyword + "COPY *.png /dest/" + ColorReset + "\n" +
		ColorKeyword + "COPY index.?s /dest/" + ColorReset + "\n" +
		"Special characters like [ and ] must be escaped using Golang rules: \n" +
		ColorKeyword + "COPY arr[[]0].txt /dest/" + ColorReset + "\n\n" +
		ColorKeyword + "Destination\n" + ColorReset +
		"If path starts with /, it is absolute. Trailing slashes are significant:\n" +
		ColorKeyword + "COPY test.txt /abs" + ColorReset + " -> /abs\n" +
		ColorKeyword + "COPY test.txt /abs/" + ColorReset + " -> /abs/test.txt\n\n" +
		"Relative paths are resolved from the WORKDIR:\n" +
		ColorKeyword + "WORKDIR /usr/src/app\n" +
		ColorKeyword + "COPY test.txt rel/" + ColorReset + " -> /usr/src/app/rel/test.txt\n\n" +
		ColorKeyword + "COPY --from\n" + ColorReset +
		"Copies from a build stage, image, or named context:\n" +
		ColorKeyword + "COPY [--from=<image|stage|context>] <src> ... <dest>" + ColorReset + "\n" +
		"For multi-stage builds, specify stage name with AS in FROM:\n" +
		ColorKeyword + "FROM alpine AS build\nCOPY . .\nRUN apk add clang\nRUN clang -o /hello hello.c\nFROM scratch\nCOPY --from=build /hello /" + ColorReset + "\n\n" +
		"Named contexts or images example:\n" +
		ColorKeyword + "COPY --from=nginx:latest /etc/nginx/nginx.conf /nginx.conf" + ColorReset + "\n\n" +
		ColorKeyword + "COPY --chown --chmod\n" + ColorReset +
		"Format:\n" +
		ColorKeyword + "COPY [--chown=<user>:<group>] [--chmod=<perms> ...] <src> ... <dest>" + ColorReset + "\n" +
		"--chown sets ownership (Linux only). Numeric IDs bypass /etc/passwd lookup. --chmod supports variable interpolation in Dockerfile 1.10+\n\n" +
		ColorKeyword + "COPY --link\n" + ColorReset +
		"Format:\n" +
		ColorKeyword + "COPY [--link[=<boolean>]] <src> ... <dest>" + ColorReset + "\n" +
		"--link copies files into an independent layer, improving cache reuse and multi-stage builds.\n\n" +
		ColorKeyword + "COPY --parents\n" + ColorReset +
		"Format:\n" +
		ColorKeyword + "COPY [--parents[=<boolean>]] <src> ... <dest>" + ColorReset + "\n" +
		"Preserves parent directories for src entries. Useful for absolute paths or multi-stage builds.\n\n" +
		ColorKeyword + "COPY --exclude\n" + ColorReset +
		"Format:\n" +
		ColorKeyword + "COPY [--exclude=<path> ...] <src> ... <dest>" + ColorReset + "\n" +
		"Excludes paths matching patterns. Multiple --exclude options can be used for a single COPY instruction.\n\n"

	info["WORKDIR"] = "\n" + ColorKeyword + "WORKDIR " + ColorCode + "/path/to/workdir" + ColorReset +
		"\n\n" +
		"The " + ColorKeyword + "WORKDIR " + ColorReset +
		"instruction sets the working directory for any " +
		ColorCode + "RUN, CMD, ENTRYPOINT, COPY " + ColorReset +
		"and " + ColorCode + "ADD " + ColorReset +
		"instructions that follow it in the Dockerfile. If the WORKDIR doesn't exist, it will be created even if it's not used in any subsequent Dockerfile instruction.\n\n" +
		"The " + ColorKeyword + "WORKDIR " + ColorReset +
		"instruction can be used multiple times in a Dockerfile. If a relative path is provided, it will be relative to the path of the previous WORKDIR instruction.\n\n" +
		ColorCode +
		"WORKDIR /a\n" +
		"WORKDIR b\n" +
		"WORKDIR c\n" +
		"RUN pwd\n" +
		ColorReset +
		"\n" +
		"The output of the final " + ColorCode + "pwd " + ColorReset +
		"command in this Dockerfile would be " +
		ColorKeyword + "/a/b/c" + ColorReset +
		".\n\n" +
		"The " + ColorKeyword + "WORKDIR " + ColorReset +
		"instruction can resolve environment variables previously set using " +
		ColorCode + "ENV" + ColorReset +
		". Only environment variables explicitly set in the Dockerfile can be used.\n\n" +
		ColorCode +
		"ENV DIRPATH=/path\n" +
		"WORKDIR $DIRPATH/$DIRNAME\n" +
		"RUN pwd\n" +
		ColorReset +
		"\n" +
		"The output of the final " + ColorCode + "pwd " + ColorReset +
		"command in this Dockerfile would be " +
		ColorKeyword + "/path/$DIRNAME" + ColorReset +
		".\n\n" +
		"If not specified, the default working directory is " +
		ColorKeyword + "/" + ColorReset +
		". In practice, if you aren't building a Dockerfile from scratch (" +
		ColorCode + "FROM scratch" + ColorReset +
		"), the WORKDIR may already be set by the base image.\n\n" +
		ColorKeyword + "Best practice\n" + ColorReset +
		"Always set " + ColorKeyword + "WORKDIR " + ColorReset +
		"explicitly to avoid unintended operations in unknown directories."
	info["ONBUILD"] = "\n" + ColorKeyword + "ONBUILD " + ColorCode + "<instruction>" + ColorReset +
		"\n\nThe " + ColorCode + "ONBUILD " + ColorReset +
		"instruction adds to the image a trigger instruction to be executed at a later time, when the image is used as the base for another build. The trigger will be executed in the context of the downstream build, as if it had been inserted immediately after the " +
		ColorCode + "FROM " + ColorReset +
		"instruction in the downstream Dockerfile.\n\nThis is useful if you are building an image which will be used as a base to build other images, for example an application build environment or a daemon which may be customized with user-specific configuration.\n\nFor example, if your image is a reusable Python application builder, it will require application source code to be added in a particular directory, and it might require a build script to be called after that. You can't just call " +
		ColorCode + "ADD " + ColorReset +
		"and " +
		ColorCode + "RUN " + ColorReset +
		"now, because you don't yet have access to the application source code, and it will be different for each application build. You could simply provide application developers with a boilerplate Dockerfile to copy-paste into their application, but that's inefficient, error-prone and difficult to update because it mixes with application-specific code.\n\nThe solution is to use " +
		ColorCode + "ONBUILD " + ColorReset +
		"to register advance instructions to run later, during the next build stage.\n\nHere's how it works:\n\n" +
		"1. When it encounters an " + ColorCode + "ONBUILD " + ColorReset + "instruction, the builder adds a trigger to the metadata of the image being built. The instruction doesn't otherwise affect the current build.\n" +
		"2. At the end of the build, a list of all triggers is stored in the image manifest, under the key " + ColorCode + "OnBuild" + ColorReset + ". They can be inspected with the " + ColorCode + "docker inspect" + ColorReset + " command.\n" +
		"3. Later the image may be used as a base for a new build, using the " + ColorCode + "FROM " + ColorReset + "instruction. As part of processing the " + ColorCode + "FROM " + ColorReset + "instruction, the downstream builder looks for " + ColorCode + "ONBUILD " + ColorReset + "triggers, and executes them in the same order they were registered. If any of the triggers fail, the " + ColorCode + "FROM " + ColorReset + "instruction is aborted which in turn causes the build to fail. If all triggers succeed, the " + ColorCode + "FROM " + ColorReset + "instruction completes and the build continues as usual.\n" +
		"4. Triggers are cleared from the final image after being executed. In other words they aren't inherited by \"grand-children\" builds.\n\nFor example you might add something like this:\n\n" +
		ColorKeyword + "ONBUILD " + ColorCode + "ADD . /app/src\n" +
		ColorKeyword + "ONBUILD " + ColorCode + "RUN /usr/local/bin/python-build --dir /app/src\n\n" +
		ColorReset +
		"Copy or mount from stage, image, or context\nAs of Dockerfile syntax 1.11, you can use " + ColorCode + "ONBUILD " + ColorReset + "with instructions that copy or mount files from other stages, images, or build contexts. For example:\n\n" +
		ColorCode + "# syntax=docker/dockerfile:1.11\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine AS baseimage\n" +
		ColorKeyword + "ONBUILD " + ColorCode + "COPY --from=build /usr/bin/app /app\n" +
		ColorKeyword + "ONBUILD " + ColorCode + "RUN --mount=from=config,target=/opt/appconfig ...\n\n" +
		ColorReset +
		"If the source of " + ColorCode + "from " + ColorReset + "is a build stage, the stage must be defined in the Dockerfile where " + ColorCode + "ONBUILD " + ColorReset + "gets triggered. If it's a named context, that context must be passed to the downstream build.\n\n" +
		ColorCode + "ONBUILD " + ColorReset +
		"limitations\n" +
		"Chaining " + ColorCode + "ONBUILD " + ColorReset + "instructions using " + ColorCode + "ONBUILD ONBUILD " + ColorReset + "isn't allowed.\n" +
		"The " + ColorCode + "ONBUILD " + ColorReset + "instruction may not trigger " + ColorCode + "FROM " + ColorReset + "or " + ColorCode + "MAINTAINER " + ColorReset + "instructions.\n\n"
	info["STOPSIGNAL"] = "\n" + ColorKeyword + "STOPSIGNAL " + ColorCode + "signal" + ColorReset +
		"\n\nThe " + ColorCode + "STOPSIGNAL " + ColorReset +
		"instruction sets the system call signal that will be sent to the container to exit. This signal can be a signal name in the format " +
		ColorCode + "SIG<NAME>" + ColorReset +
		", for instance " +
		ColorCode + "SIGKILL" + ColorReset +
		", or an unsigned number that matches a position in the kernel's syscall table, for instance " +
		ColorCode + "9" + ColorReset +
		". The default is " +
		ColorCode + "SIGTERM" + ColorReset +
		" if not defined.\n\nThe image's default stopsignal can be overridden per container, using the " +
		ColorCode + "--stop-signal " + ColorReset +
		"flag on " +
		ColorCode + "docker run " + ColorReset +
		"and " +
		ColorCode + "docker create" + ColorReset + ".\n\n"
	info["HEALTHCHECK"] = "\n" + ColorKeyword + "HEALTHCHECK " + ColorCode + "[OPTIONS] CMD command" + ColorReset +
		" (check container health by running a command inside the container)\n" +
		ColorKeyword + "HEALTHCHECK " + ColorCode + "NONE" + ColorReset +
		" (disable any healthcheck inherited from the base image)\n\n" +
		"The " + ColorCode + "HEALTHCHECK " + ColorReset +
		"instruction tells Docker how to test a container to check that it's still working. This can detect cases such as a web server stuck in an infinite loop and unable to handle new connections, even though the server process is still running.\n\n" +
		"When a container has a " + ColorCode + "healthcheck " + ColorReset + "specified, it has a health status in addition to its normal status. This status is initially " +
		ColorCode + "starting" + ColorReset + ". Whenever a health check passes, it becomes " +
		ColorCode + "healthy" + ColorReset + " (whatever state it was previously in). After a certain number of consecutive failures, it becomes " +
		ColorCode + "unhealthy" + ColorReset + ".\n\n" +
		"The options that can appear before " + ColorCode + "CMD " + ColorReset + "are:\n\n" +
		ColorCode + "--interval=DURATION" + ColorReset + " (default: 30s)\n" +
		ColorCode + "--timeout=DURATION" + ColorReset + " (default: 30s)\n" +
		ColorCode + "--start-period=DURATION" + ColorReset + " (default: 0s)\n" +
		ColorCode + "--start-interval=DURATION" + ColorReset + " (default: 5s)\n" +
		ColorCode + "--retries=N" + ColorReset + " (default: 3)\n\n" +
		"The health check will first run " + ColorCode + "interval" + ColorReset + " seconds after the container is started, and then again interval seconds after each previous check completes.\n\n" +
		"If a single run of the check takes longer than " + ColorCode + "timeout" + ColorReset + " seconds then the check is considered to have failed. The process performing the check is abruptly stopped with a " +
		ColorCode + "SIGKILL" + ColorReset + ".\n\n" +
		"It takes " + ColorCode + "retries" + ColorReset + " consecutive failures of the health check for the container to be considered unhealthy.\n\n" +
		"Start period provides initialization time for containers that need time to bootstrap. Probe failure during that period will not be counted towards the maximum number of retries. However, if a health check succeeds during the start period, the container is considered started and all consecutive failures will be counted towards the maximum number of retries.\n\n" +
		"Start interval is the time between health checks during the start period. This option requires Docker Engine version 25.0 or later.\n\n" +
		"There can only be one " + ColorCode + "HEALTHCHECK " + ColorReset + "instruction in a Dockerfile. If you list more than one then only the last HEALTHCHECK will take effect.\n\n" +
		"The command after the " + ColorCode + "CMD " + ColorReset + "keyword can be either a shell command (e.g. " +
		ColorKeyword + "HEALTHCHECK CMD /bin/check-running" + ColorReset + ") or an exec array (as with other Dockerfile commands; see e.g. ENTRYPOINT for details).\n\n" +
		"The command's exit status indicates the health status of the container. The possible values are:\n\n" +
		ColorCode + "0" + ColorReset + ": success - the container is healthy and ready for use\n" +
		ColorCode + "1" + ColorReset + ": unhealthy - the container isn't working correctly\n" +
		ColorCode + "2" + ColorReset + ": reserved - don't use this exit code\n\n" +
		"For example, to check every five minutes or so that a web-server is able to serve the site's main page within three seconds:\n\n" +
		ColorKeyword + "HEALTHCHECK --interval=5m --timeout=3s \\\n  CMD curl -f http://localhost/ || exit 1\n\n" +
		"To help debug failing probes, any output text (UTF-8 encoded) that the command writes on stdout or stderr will be stored in the health status and can be queried with " +
		ColorCode + "docker inspect" + ColorReset + ". Such output should be kept short (only the first 4096 bytes are stored currently).\n\n" +
		"When the health status of a container changes, a " + ColorCode + "health_status " + ColorReset + "event is generated with the new status.\n\n"
	info["SHELL"] = "\n" + ColorKeyword + "SHELL " + ColorCode + "[\"executable\", \"parameters\"]" + ColorReset +
		"\n\nThe " + ColorCode + "SHELL " + ColorReset +
		"instruction allows the default shell used for the shell form of commands to be overridden. The default shell on Linux is " +
		ColorCode + "[\"/bin/sh\", \"-c\"]" + ColorReset +
		", and on Windows is " +
		ColorCode + "[\"cmd\", \"/S\", \"/C\"]" + ColorReset +
		". The " + ColorCode + "SHELL " + ColorReset +
		"instruction must be written in JSON form in a Dockerfile.\n\n" +
		"The " + ColorCode + "SHELL " + ColorReset +
		"instruction is particularly useful on Windows where there are two commonly used and quite different native shells: " +
		ColorCode + "cmd " + ColorReset + "and " +
		ColorCode + "powershell" + ColorReset +
		", as well as alternate shells available including sh.\n\n" +
		"The " + ColorCode + "SHELL " + ColorReset +
		"instruction can appear multiple times. Each " + ColorCode + "SHELL " + ColorReset +
		"instruction overrides all previous " + ColorCode + "SHELL " + ColorReset +
		"instructions, and affects all subsequent instructions. For example:\n\n" +
		ColorKeyword + "FROM " + ColorCode + "microsoft/windowsservercore\n\n" +
		ColorReset + "# Executed as cmd /S /C echo default\n" +
		ColorKeyword + "RUN " + ColorCode + "echo default\n\n" +
		ColorReset + "# Executed as cmd /S /C powershell -command Write-Host default\n" +
		ColorKeyword + "RUN " + ColorCode + "powershell -command Write-Host default\n\n" +
		ColorReset + "# Executed as powershell -command Write-Host hello\n" +
		ColorKeyword + "SHELL " + ColorCode + "[\"powershell\", \"-command\"]\n" +
		ColorKeyword + "RUN " + ColorCode + "Write-Host hello\n\n" +
		ColorReset + "# Executed as cmd /S /C echo hello\n" +
		ColorKeyword + "SHELL " + ColorCode + "[\"cmd\", \"/S\", \"/C\"]\n" +
		ColorKeyword + "RUN " + ColorCode + "echo hello\n\n" +
		"The following instructions can be affected by the " + ColorCode + "SHELL " + ColorReset +
		"instruction when the shell form of them is used in a Dockerfile: " +
		ColorCode + "RUN, CMD, ENTRYPOINT" + ColorReset + ".\n\n" +
		"The following example is a common pattern found on Windows which can be streamlined by using the " + ColorCode + "SHELL " + ColorReset +
		"instruction:\n\n" +
		ColorKeyword + "RUN " + ColorCode + "powershell -command Execute-MyCmdlet -param1 \"c:\\foo.txt\"\n\n" +
		"The command invoked by the builder will be:\n\n" +
		ColorCode + "cmd /S /C powershell -command Execute-MyCmdlet -param1 \"c:\\foo.txt\"\n\n" +
		"This is inefficient for two reasons. First, there is an unnecessary cmd.exe command processor (aka shell) being invoked. Second, each RUN instruction in the shell form requires an extra powershell -command prefixing the command.\n\n" +
		"To make this more efficient, one of two mechanisms can be employed. One is to use the JSON form of the RUN command such as:\n\n" +
		ColorKeyword + "RUN " + ColorCode + "[\"powershell\", \"-command\", \"Execute-MyCmdlet\", \"-param1 \\\"c:\\\\foo.txt\\\"\"]\n\n" +
		"While the JSON form is unambiguous and does not use the unnecessary cmd.exe, it does require more verbosity through double-quoting and escaping. The alternate mechanism is to use the " + ColorCode + "SHELL " + ColorReset +
		"instruction and the shell form, making a more natural syntax for Windows users, especially when combined with the escape parser directive:\n\n" +
		ColorCode + "# escape=`\n" +
		ColorKeyword + "FROM " + ColorCode + "microsoft/nanoserver\n" +
		ColorKeyword + "SHELL " + ColorCode + "[\"powershell\",\"-command\"]\n" +
		ColorKeyword + "RUN " + ColorCode + "New-Item -ItemType Directory C:\\Example\n" +
		ColorKeyword + "ADD " + ColorCode + "Execute-MyCmdlet.ps1 c:\\example\\\n" +
		ColorKeyword + "RUN " + ColorCode + "c:\\example\\Execute-MyCmdlet -sample 'hello world'\n\n" +
		ColorReset + "Resulting in:\n\n" +
		ColorCode + "PS E:\\myproject> docker build -t shell .\n" +
		"Sending build context to Docker daemon 4.096 kB\n" +
		"Step 1/5 : FROM microsoft/nanoserver\n" +
		" ---> 22738ff49c6d\n" +
		"Step 2/5 : SHELL powershell -command\n" +
		" ---> Running in 6fcdb6855ae2\n" +
		" ---> 6331462d4300\n" +
		"Removing intermediate container 6fcdb6855ae2\n" +
		"Step 3/5 : RUN New-Item -ItemType Directory C:\\Example\n" +
		" ---> Running in d0eef8386e97\n\n" +
		"    Directory: C:\\\n\n" +
		"Mode         LastWriteTime              Length Name\n" +
		"----         -------------              ------ ----\n" +
		"d-----       10/28/2016  11:26 AM              Example\n\n" +
		" ---> 3f2fbf1395d9\n" +
		"Removing intermediate container d0eef8386e97\n" +
		"Step 4/5 : ADD Execute-MyCmdlet.ps1 c:\\example\\\n" +
		" ---> a955b2621c31\n" +
		"Removing intermediate container b825593d39fc\n" +
		"Step 5/5 : RUN c:\\example\\Execute-MyCmdlet 'hello world'\n" +
		" ---> Running in be6d8e63fe75\n" +
		"hello world\n" +
		" ---> 8e559e9bf424\n" +
		"Removing intermediate container be6d8e63fe75\n" +
		"Successfully built 8e559e9bf424\n" +
		"PS E:\\myproject>\n\n" +
		"The " + ColorCode + "SHELL " + ColorReset + "instruction could also be used to modify the way in which a shell operates. For example, using " +
		ColorCode + "SHELL cmd /S /C /V:ON|OFF" + ColorReset +
		" on Windows, delayed environment variable expansion semantics could be modified.\n\n" +
		"The " + ColorCode + "SHELL " + ColorReset + "instruction can also be used on Linux should an alternate shell be required such as " +
		ColorCode + "zsh, csh, tcsh" + ColorReset +
		" and others.\n\n" +
		ColorKeyword + "Here-Documents\n" + ColorReset +
		"Here-documents allow redirection of subsequent Dockerfile lines to the input of " + ColorCode + "RUN " + ColorReset + "or " + ColorCode + "COPY " + ColorReset + "commands. If such command contains a here-document the Dockerfile considers the next lines until the line only containing a here-doc delimiter as part of the same command.\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "debian\n" +
		ColorKeyword + "RUN " + ColorCode + "<<EOT bash\n" +
		"  set -ex\n" +
		"  apt-get update\n" +
		"  apt-get install -y vim\n" +
		"EOT\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "debian\n" +
		ColorKeyword + "RUN " + ColorCode + "<<EOT\n" +
		"  mkdir -p foo/bar\n" +
		"EOT\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "python:3.6\n" +
		ColorKeyword + "RUN " + ColorCode + "<<EOT\n" +
		"#!/usr/bin/env python\n" +
		"print(\"hello world\")\n" +
		"EOT\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine\n" +
		ColorKeyword + "RUN " + ColorCode + "<<FILE1 cat > file1 && <<FILE2 cat > file2\n" +
		"I am\nfirst\nFILE1\nI am\nsecond\nFILE2\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine\n" +
		ColorKeyword + "COPY " + ColorCode + "<<EOF greeting.txt\n" +
		"hello world\n" +
		"EOF\n\n" +
		ColorReset + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine\n" +
		ColorKeyword + "ARG " + ColorCode + "FOO=bar\n" +
		ColorKeyword + "COPY " + ColorCode + "<<-EOT /script.sh\n" +
		"  echo \"hello ${FOO}\"\n" +
		"EOT\n" +
		ColorKeyword + "ENTRYPOINT " + ColorCode + "ash /script.sh\n\n" +
		ColorReset + "In this case, file script prints " + ColorCode + "\"hello bar\"" + ColorReset + ", because the variable is expanded when the COPY instruction gets executed.\n\n" +
		ColorCode + "docker build -t heredoc .\n" +
		ColorCode + "docker run heredoc\n" +
		"hello bar\n\n" +
		"Quoted here-documents prevent expansion at build-time:\n\n" +
		ColorKeyword + "# syntax=docker/dockerfile:1\n" +
		ColorKeyword + "FROM " + ColorCode + "alpine\n" +
		ColorKeyword + "ARG " + ColorCode + "FOO=bar\n" +
		ColorKeyword + "COPY " + ColorCode + "<<-\"EOT\" /script.sh\n" +
		"  echo \"hello ${FOO}\"\n" +
		"EOT\n" +
		ColorKeyword + "ENTRYPOINT " + ColorCode + "ash /script.sh\n\n" +
		ColorReset + "Variable is interpreted at runtime when the script is invoked:\n\n" +
		ColorCode + "docker build -t heredoc .\n" +
		ColorCode + "docker run -e FOO=world heredoc\n" +
		"hello world\n\n"
	info["RUN"] = "\n" + ColorKeyword + "RUN " + ColorReset + "instruction will execute any commands to create a new layer on top of the current image. The added layer is used in the next step in the Dockerfile. " +
		ColorKeyword + "RUN " + ColorReset + "has two forms:\n\n" +
		ColorComment + "# Shell form:" + ColorReset + "\n" +
		ColorKeyword + "RUN " + ColorCode + "[OPTIONS] <command> ..." + ColorReset + "\n" +
		ColorComment + "# Exec form:" + ColorReset + "\n" +
		ColorKeyword + "RUN " + ColorCode + "[OPTIONS] [ \"<command>\", ... ]" + ColorReset + "\n\n" +
		"The shell form is most commonly used, and lets you break up longer instructions into multiple lines, either using newline escapes, or with heredocs:\n\n" +
		ColorKeyword + "RUN " + ColorCode + "<<EOF\napt-get update\napt-get install -y curl\nEOF" + ColorReset + "\n\n" +
		"The available " + ColorCode + "[OPTIONS]" + ColorReset + " for the " + ColorKeyword + "RUN " + ColorReset + "instruction are:\n\n" +
		ColorCode + "Option\tMinimum Dockerfile version" + ColorReset + "\n" +
		"--device\t1.14-labs\n" +
		"--mount\t1.2\n" +
		"--network\t1.3\n" +
		"--security\t1.20\n\n" +
		"Cache invalidation for " + ColorKeyword + "RUN " + ColorReset + "instructions\n" +
		"The cache for " + ColorKeyword + "RUN " + ColorReset + "instructions isn't invalidated automatically during the next build. The cache for an instruction like " +
		ColorCode + "RUN apt-get dist-upgrade -y" + ColorReset + " will be reused during the next build. The cache for " + ColorKeyword + "RUN " + ColorReset + "instructions can be invalidated by using the " +
		ColorCode + "--no-cache " + ColorReset + "flag, for example " + ColorCode + "docker build --no-cache" + ColorReset + ".\n\n" +
		"See the Dockerfile Best Practices guide for more information.\n\n" +
		"The cache for " + ColorKeyword + "RUN " + ColorReset + "instructions can be invalidated by " + ColorKeyword + "ADD " + ColorReset + "and " + ColorKeyword + "COPY " + ColorReset + "instructions.\n\n" +
		ColorKeyword + "RUN --device" + ColorReset + "\n" +
		ColorComment + "Note\n" + ColorReset +
		"Not yet available in stable syntax, use " + ColorCode + "docker/dockerfile:1-labs" + ColorReset + " version. It also needs BuildKit 0.20.0 or later.\n\n" +
		ColorKeyword + "RUN --device=name,[required]" + ColorReset + "\n" +
		ColorKeyword + "RUN --device " + ColorReset + "allows build to request CDI devices to be available to the build step.\n\n" +
		ColorComment + "Warning\n" + ColorReset +
		"The use of " + ColorCode + "--device " + ColorReset + "is protected by the device entitlement, which needs to be enabled when starting the buildkitd daemon with " +
		ColorCode + "--allow-insecure-entitlement device" + ColorReset + " flag or in buildkitd config, and for a build request with " +
		ColorCode + "--allow device" + ColorReset + " flag.\n\n" +
		"The device name is provided by the CDI specification registered in BuildKit.\n\n" +
		"Example: multiple devices registered in CDI spec:\n\n" +
		ColorCode + "cdiVersion: \"0.6.0\"\nkind: \"vendor1.com/device\"\ndevices:\n  - name: foo\n    containerEdits:\n      env:\n        - FOO=injected\n  - name: bar\n    annotations:\n      org.mobyproject.buildkit.device.class: class1\n    containerEdits:\n      env:\n        - BAR=injected\n  - name: baz\n    annotations:\n      org.mobyproject.buildkit.device.class: class1\n    containerEdits:\n      env:\n        - BAZ=injected\n  - name: qux\n    annotations:\n      org.mobyproject.buildkit.device.class: class2\n    containerEdits:\n      env:\n        - QUX=injected\nannotations:\n  org.mobyproject.buildkit.device.autoallow: true" + ColorReset + "\n\n" +
		"Device name format supports various patterns: vendor1.com/device, vendor1.com/device=foo, vendor1.com/device=*, class1.\n\n" +
		ColorComment + "Note\n" + ColorReset +
		"Annotations supported by CDI spec since 0.6.0. You can also set org.mobyproject.buildkit.device.autoallow annotation for all devices or specific device.\n\n" +
		"Example: CUDA-Powered LLaMA Inference:\n\n" +
		ColorCode + "# syntax=docker/dockerfile:1-labs\nFROM scratch AS model\nADD https://huggingface.co/.../model.gguf /model.gguf\n\nFROM scratch AS prompt\nCOPY <<EOF prompt.txt\nQ: Generate ...\nEOF\n\nFROM ghcr.io/ggml-org/llama.cpp:full-cuda-b5124\nRUN --device=nvidia.com/gpu=all \\\n    --mount=from=model,target=/models \\\n    --mount=from=prompt,target=/tmp \\\n    ./llama-cli -m /models/model.gguf -no-cnv -ngl 99 -f /tmp/prompt.txt" + ColorReset + "\n\n" +
		ColorKeyword + "RUN --mount" + ColorReset + "\n" +
		ColorKeyword + "RUN --mount=[type=TYPE][,option=<value>[,option=<value>]...]" + ColorReset + " allows you to create filesystem mounts that the build can access. This can be used to:\n\n" +
		"- Create bind mount to host filesystem or other build stages\n- Access build secrets or ssh-agent sockets\n- Use persistent package cache\n\n" +
		"Supported mount types: bind, cache, tmpfs, secret, ssh.\n\n" +
		ColorKeyword + "RUN --mount=type=bind" + ColorReset + " ... (details)\n" +
		ColorKeyword + "RUN --mount=type=cache" + ColorReset + " ... (details)\n" +
		ColorKeyword + "RUN --mount=type=tmpfs" + ColorReset + " ... (details)\n" +
		ColorKeyword + "RUN --mount=type=secret" + ColorReset + " ... (details)\n" +
		ColorKeyword + "RUN --mount=type=ssh" + ColorReset + " ... (details)\n\n" +
		ColorKeyword + "RUN --network" + ColorReset + " ... (details)\n" +
		ColorKeyword + "RUN --security" + ColorReset + " ... (details)\n"
	info["FROM"] = "\n" + ColorKeyword + "FROM " + ColorCode + "[--platform=<platform>] <image> [AS <name>]" + ColorReset + "\n" +
		"Or\n\n" +
		ColorKeyword + "FROM " + ColorCode + "[--platform=<platform>] <image>[:<tag>] [AS <name>]" + ColorReset + "\n" +
		"Or\n\n" +
		ColorKeyword + "FROM " + ColorCode + "[--platform=<platform>] <image>[@<digest>] [AS <name>]" + ColorReset + "\n\n" +
		"The " + ColorKeyword + "FROM " + ColorReset + "instruction initializes a new build stage and sets the base image for subsequent instructions. " +
		"As such, a valid Dockerfile must start with a " + ColorKeyword + "FROM " + ColorReset + "instruction. The image can be any valid image.\n\n" +
		ColorKeyword + "ARG " + ColorReset + "is the only instruction that may precede " + ColorKeyword + "FROM " + ColorReset + "in the Dockerfile. See " + ColorCode + "Understand how ARG and FROM interact" + ColorReset + ".\n\n" +
		ColorKeyword + "FROM " + ColorReset + "can appear multiple times within a single Dockerfile to create multiple images or use one build stage as a dependency for another. " +
		"Simply make a note of the last image ID output by the commit before each new " + ColorKeyword + "FROM " + ColorReset + "instruction. Each " + ColorKeyword + "FROM " + ColorReset + "instruction clears any state created by previous instructions.\n\n" +
		"Optionally a name can be given to a new build stage by adding " + ColorCode + "AS <name> " + ColorReset + "to the " + ColorKeyword + "FROM " + ColorReset + "instruction. The name can be used in subsequent " +
		ColorKeyword + "FROM <name>, COPY --from=<name>, " + ColorKeyword + "RUN --mount=type=bind,from=<name>" + ColorReset + " instructions to refer to the image built in this stage.\n\n" +
		"The tag or digest values are optional. If you omit either of them, the builder assumes a latest tag by default. The builder returns an error if it can't find the tag value.\n\n" +
		"The optional " + ColorCode + "--platform" + ColorReset + " flag can be used to specify the platform of the image in case " + ColorKeyword + "FROM " + ColorReset + "references a multi-platform image. " +
		"For example, linux/amd64, linux/arm64, or windows/amd64. By default, the target platform of the build request is used. " +
		"Global build arguments can be used in the value of this flag, for example automatic platform ARGs allow you to force a stage to native build platform (" +
		ColorCode + "--platform=$BUILDPLATFORM" + ColorReset + "), and use it to cross-compile to the target platform inside the stage.\n\n" +
		ColorCode + "Understand how ARG and FROM interact" + ColorReset + "\n" +
		ColorKeyword + "FROM " + ColorReset + "instructions support variables that are declared by any " + ColorKeyword + "ARG " + ColorReset + "instructions that occur before the first " + ColorKeyword + "FROM" + ColorReset + ".\n\n" +
		ColorKeyword + "ARG " + ColorCode + "CODE_VERSION=latest" + ColorReset + "\n" +
		ColorKeyword + "FROM " + ColorCode + "base:${CODE_VERSION}" + ColorReset + "\n" +
		ColorKeyword + "CMD " + ColorCode + "/code/run-app" + ColorReset + "\n\n" +
		ColorKeyword + "FROM " + ColorCode + "extras:${CODE_VERSION}" + ColorReset + "\n" +
		ColorKeyword + "CMD " + ColorCode + "/code/run-extras" + ColorReset + "\n\n" +
		"An " + ColorKeyword + "ARG " + ColorReset + "declared before a " + ColorKeyword + "FROM " + ColorReset + "is outside of a build stage, so it can't be used in any instruction after a " +
		ColorKeyword + "FROM" + ColorReset + ". To use the default value of an " + ColorKeyword + "ARG " + ColorReset + "declared before the first " + ColorKeyword + "FROM " + ColorReset + "use an " +
		ColorKeyword + "ARG " + ColorReset + "instruction without a value inside of a build stage:\n\n" +
		ColorKeyword + "ARG " + ColorCode + "VERSION=latest" + ColorReset + "\n" +
		ColorKeyword + "FROM " + ColorCode + "busybox:$VERSION" + ColorReset + "\n" +
		ColorKeyword + "ARG " + ColorReset + "VERSION\n" +
		ColorKeyword + "RUN " + ColorCode + "echo $VERSION > image_version" + ColorReset + "\n\n"
	info["ARG"] = "\n" + ColorKeyword + "ARG " + ColorCode + "<name>[=<default value>] [<name>[=<default value>]...]" + ColorReset +
		"\n\nThe " + ColorCode + "ARG " + ColorReset +
		"instruction defines a variable that users can pass at build-time to the builder with the " + ColorCode + "docker build" + ColorReset +
		" command using the " + ColorCode + "--build-arg <varname>=<value>" + ColorReset + " flag.\n\n" +
		ColorCode + "Warning" + ColorReset + "\nIt isn't recommended to use build arguments for passing secrets such as user credentials, API tokens, etc. Build arguments are visible in the " + ColorCode + "docker history" + ColorReset +
		" command and in max mode provenance attestations, which are attached to the image by default if you use the Buildx GitHub Actions and your GitHub repository is public.\n\n" +
		"Refer to the " + ColorCode + "RUN --mount=type=secret" + ColorReset + " section to learn about secure ways to use secrets when building images.\n\n" +
		"A Dockerfile may include one or more " + ColorKeyword + "ARG " + ColorReset + "instructions. For example, the following is a valid Dockerfile:\n\n" +
		ColorKeyword + "FROM busybox\n" +
		ColorKeyword + "ARG " + ColorCode + "user1\n" +
		ColorKeyword + "ARG " + ColorCode + "buildno\n" +
		"# ...\n\n" +
		ColorCode + "Default values" + ColorReset + "\nAn " + ColorKeyword + "ARG " + ColorReset + "instruction can optionally include a default value:\n\n" +
		ColorKeyword + "FROM busybox\n" +
		ColorKeyword + "ARG " + ColorCode + "user1=someuser\n" +
		ColorKeyword + "ARG " + ColorCode + "buildno=1\n" +
		"# ...\n\n" +
		"If an " + ColorKeyword + "ARG " + ColorReset + "instruction has a default value and if there is no value passed at build-time, the builder uses the default.\n\n" +
		ColorCode + "Scope" + ColorReset + "\nAn " + ColorKeyword + "ARG " + ColorReset + "variable comes into effect from the line on which it is declared in the Dockerfile. For example, consider this Dockerfile:\n\n" +
		ColorKeyword + "FROM busybox\n" +
		ColorKeyword + "USER ${username:-some_user}\n" +
		ColorKeyword + "ARG " + ColorCode + "username\n" +
		ColorKeyword + "USER $username\n" +
		"# ...\n\n" +
		"A user builds this file by calling:\n\n " +
		ColorCode + "docker build --build-arg username=what_user ." + ColorReset + "\n\n" +
		"The " + ColorKeyword + "USER " + ColorReset + "instruction on line 2 evaluates to the " + ColorCode + "some_user" + ColorReset +
		" fallback, because the username variable is not yet declared. The username variable is declared on line 3, and available for reference in Dockerfile instructions from that point onwards. The " + ColorKeyword + "USER " + ColorReset +
		"instruction on line 4 evaluates to " + ColorCode + "what_user" + ColorReset + ", since at that point the username argument has a value of what_user which was passed on the command line. Prior to its definition by an " +
		ColorKeyword + "ARG " + ColorReset + "instruction, any use of a variable results in an empty string.\n\n" +
		"An " + ColorKeyword + "ARG " + ColorReset + "variable declared within a build stage is automatically inherited by other stages based on that stage. Unrelated build stages do not have access to the variable. To use an argument in multiple distinct stages, each stage must include the " + ColorKeyword + "ARG " + ColorReset +
		"instruction, or they must both be based on a shared base stage in the same Dockerfile where the variable is declared.\n\n" +
		ColorCode + "Using ARG variables" + ColorReset + "\nYou can use an " + ColorKeyword + "ARG " + ColorReset + "or an " + ColorKeyword + "ENV " + ColorReset + "instruction to specify variables that are available to the " + ColorKeyword + "RUN " + ColorReset +
		"instruction. Environment variables defined using the " + ColorKeyword + "ENV " + ColorReset + "instruction always override an " + ColorKeyword + "ARG " + ColorReset + "instruction of the same name. Consider this Dockerfile with an " +
		ColorKeyword + "ENV " + ColorReset + "and " + ColorKeyword + "ARG " + ColorReset + "instruction:\n\n" +
		ColorKeyword + "FROM ubuntu\n" +
		ColorKeyword + "ARG " + ColorCode + "CONT_IMG_VER\n" +
		ColorKeyword + "ENV CONT_IMG_VER=v1.0.0\n" +
		ColorKeyword + "RUN echo $CONT_IMG_VER\n\n" +
		"Then, assume this image is built with this command:\n\n " +
		ColorCode + "docker build --build-arg CONT_IMG_VER=v2.0.1 ." + ColorReset + "\n\n" +
		"In this case, the " + ColorKeyword + "RUN " + ColorReset + "instruction uses " + ColorCode + "v1.0.0" + ColorReset + " instead of the " + ColorKeyword + "ARG " + ColorReset +
		"setting passed by the user: " + ColorCode + "v2.0.1" + ColorReset + ". This behavior is similar to a shell script where a locally scoped variable overrides the variables passed as arguments or inherited from environment, from its point of definition.\n\n" +
		"Unlike an " + ColorKeyword + "ARG " + ColorReset + "instruction, " + ColorKeyword + "ENV " + ColorReset + "values are always persisted in the built image. Consider a docker build without the " + ColorCode + "--build-arg" + ColorReset + " flag:\n\n " +
		ColorCode + "docker build ." + ColorReset + "\n\n" +
		"Using this Dockerfile example, CONT_IMG_VER is still persisted in the image but its value would be " + ColorCode + "v1.0.0" + ColorReset + " as it is the default set in line 3 by the " + ColorKeyword + "ENV " + ColorReset + "instruction.\n\n" +
		"The variable expansion technique in this example allows you to pass arguments from the command line and persist them in the final image by leveraging the " + ColorKeyword + "ENV " + ColorReset + "instruction. Variable expansion is only supported for a limited set of Dockerfile instructions.\n\n" +
		ColorCode + "Predefined ARGs" + ColorReset + "\nDocker has a set of predefined " + ColorKeyword + "ARG " + ColorReset + "variables that you can use without a corresponding ARG instruction in the Dockerfile:\n\n" +
		ColorCode + "HTTP_PROXY\nhttp_proxy\nHTTPS_PROXY\nhttps_proxy\nFTP_PROXY\nftp_proxy\nNO_PROXY\nno_proxy\nALL_PROXY\nall_proxy\n\n" +
		"To use these, pass them on the command line using the " + ColorCode + "--build-arg" + ColorReset + " flag, for example:\n\n " +
		ColorCode + "docker build --build-arg HTTPS_PROXY=https://my-proxy.example.com ." + ColorReset + "\n\n" +
		"By default, these predefined variables are excluded from the output of " + ColorKeyword + "docker history" + ColorReset + ". Excluding them reduces the risk of accidentally leaking sensitive authentication information in an " + ColorCode + "HTTP_PROXY " + ColorReset + "variable.\n\n" +
		"Automatic platform " + ColorKeyword + "ARGs " + ColorReset + "in the global scope\nThis feature is only available when using the BuildKit backend.\n\n" +
		"BuildKit supports a predefined set of " + ColorKeyword + "ARG " + ColorReset + "variables with information on the platform of the node performing the build (build platform) and on the platform of the resulting image (target platform). The target platform can be specified with the " + ColorCode + "--platform" + ColorReset + " flag on " + ColorCode + "docker build" + ColorReset + ".\n\n" +
		"Example variables:\nTARGETPLATFORM, TARGETOS, TARGETARCH, TARGETVARIANT, BUILDPLATFORM, BUILDOS, BUILDARCH, BUILDVARIANT\n\n" +
		"These arguments are defined in the global scope so are not automatically available inside build stages or for your " + ColorKeyword + "RUN " + ColorReset + "commands. To expose one of these arguments inside the build stage redefine it without value.\n\n" +
		"Impact on build caching\n" +
		"ARG variables are not persisted into the built image as " + ColorKeyword + "ENV " + ColorReset + "variables are. However, ARG variables do impact the build cache in similar ways. If a Dockerfile defines an ARG variable whose value is different from a previous build, then a 'cache miss' occurs upon its first usage, not its definition.\n\n"

}

func (e DockerfileEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		e.applySyntaxHighlighting(v)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		e.applySyntaxHighlighting(v)
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
		e.applySyntaxHighlighting(v)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
		e.applySyntaxHighlighting(v)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		v.EditNewLine()
		e.applySyntaxHighlighting(v)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
	case key == gocui.KeyHome:
		_, cy := v.Cursor()
		v.SetCursor(0, cy)
	case key == gocui.KeyEnd:
		_, cy := v.Cursor()
		line, _ := v.Line(cy)
		v.SetCursor(len(line), cy)
	}
}

func (e DockerfileEditor) applySyntaxHighlighting(v *gocui.View) {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	content := stripAnsiCodes(v.Buffer())

	highlighted := highlightDockerfile(content)

	v.Clear()
	fmt.Fprint(v, highlighted)

	v.SetCursor(cx, cy)
	v.SetOrigin(ox, oy)
}

func stripAnsiCodes(text string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(text, "")
}

func highlightDockerfile(content string) string {
	lines := strings.Split(content, "\n")
	var highlightedLines []string

	for _, line := range lines {
		highlightedLines = append(highlightedLines, highlightLine(line))
	}
	return strings.Join(highlightedLines, "\n")
}

func highlightLine(line string) string {
	trimmedLine := strings.TrimSpace(line)

	if strings.HasPrefix(trimmedLine, "#") {
		return ColorComment + line + ColorReset
	}

	for _, keyword := range dockerfileKeywords {
		if strings.HasPrefix(strings.ToUpper(trimmedLine), keyword) {
			keyWordEnd := len(keyword)
			if len(trimmedLine) >= keyWordEnd {
				leadingSpaces := len(line) - len(trimmedLine)
				actualKeyword := line[leadingSpaces : leadingSpaces+keyWordEnd]
				rest := ""
				if len(line) > leadingSpaces+keyWordEnd {
					rest = line[leadingSpaces+keyWordEnd:]
				}

				highlighted := line[:leadingSpaces] + ColorKeyword + actualKeyword + ColorReset

				rest = highlightStrings(rest)

				return highlighted + rest
			}
		}
	}
	return highlightStrings(line)
}

func highlightStrings(text string) string {
	stringRegex := regexp.MustCompile(`("[^"]*"|'[^']*')`)

	return stringRegex.ReplaceAllStringFunc(text, func(match string) string {
		return ColorString + match + ColorReset
	})
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
	}
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true
		v.Editor = DockerfileEditor{}

		if _, err := g.SetCurrentView("body"); err != nil {
			return err
		}
		if fileExists(file) {
			b, err := os.ReadFile(file)
			if err != nil {
				panic(err)
			}

			content := string(b)
			highlighted := highlightDockerfile(content)
			fmt.Fprint(v, highlighted)
		}

	}
	if v, err := g.SetView("help", 2*maxX/3+1, 3, maxX-1, 3*maxY/4, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "\033[32mADD")
		fmt.Fprintln(v, "ARG")
		fmt.Fprintln(v, "CMD")
		fmt.Fprintln(v, "COPY")
		fmt.Fprintln(v, "ENTRYPOINT")
		fmt.Fprintln(v, "ENV")
		fmt.Fprintln(v, "EXPOSE")
		fmt.Fprintln(v, "FROM")
		fmt.Fprintln(v, "HEALTHCHECK")
		fmt.Fprintln(v, "LABEL")
		fmt.Fprintln(v, "MAINTAINER")
		fmt.Fprintln(v, "ONBUILD")
		fmt.Fprintln(v, "RUN")
		fmt.Fprintln(v, "SHELL")
		fmt.Fprintln(v, "STOPSIGNAL")
		fmt.Fprintln(v, "USER")
		fmt.Fprintln(v, "VOLUME")
		fmt.Fprintln(v, "WORKDIR\033[0m")
	}
	if v, err := g.SetView("command", 2*maxX/3+1, 3*maxY/4+1, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Commands for using the editor")
		fmt.Fprintln(v, "Ctrl-C: Exit")
		fmt.Fprintln(v, "Ctrl-S: Save")
		fmt.Fprintln(v, "Ctrl-O: Toggle Overwrite")
		fmt.Fprintln(v, "Ctrl-H: Toggle Help")
	}
	return nil
}

func newFile(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, "Dockerfile")
	if err != nil {
		return err
	}

	content := stripAnsiCodes(v.Buffer())

	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return err
	}

	if err := os.Rename(f.Name(), file); err != nil {
		os.Remove(f.Name())
		return err
	}
	return nil
}
func nothing(g *gocui.Gui, v *gocui.View) error {
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func overwrite(g *gocui.Gui, v *gocui.View) error {
	v.Overwrite = !v.Overwrite
	return nil
}
func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return err
	}
	maxX, maxY := g.Size()
	if v, err := g.SetView("information", 2*maxX/3+1, 3, maxX-1, 3*maxY/4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		fmt.Fprintln(v, ColorKeyword+l+ColorReset)
		v.Wrap = true
		g.Cursor = true
		v.Editable = true
		for _, keyword := range dockerfileKeywords {
			if l == keyword {
				fmt.Fprintln(v, info[keyword])
			}
		}
		if _, err := g.SetCurrentView("information"); err != nil {
			return err
		}
	}
	return nil
}
func deleteInformationView(g *gocui.Gui, v *gocui.View) error {
	if !viewExist(g, "information") {
		return nil
	}
	if err := g.DeleteView("information"); err != nil {
		return err
	}
	if err := bodyToHelpView(g, v); err != nil {
		return err
	}
	return nil
}
func viewExist(g *gocui.Gui, s string) bool {
	if _, err := g.View(s); err != nil {
		return false
	}
	return true
}
func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}
func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, deleteInformationView); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlO, gocui.ModNone, overwrite); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlN, gocui.ModNone, newFile); err != nil {
		return err
	}
	if err := g.SetKeybinding("body", gocui.KeyCtrlS, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return saveView(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("savename", gocui.KeyEnter, gocui.ModNone, saveDeleteView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.MouseWheelDown, gocui.ModNone, nothing); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.MouseWheelUp, gocui.ModNone, nothing); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyCtrlH, gocui.ModNone, helpToBodyView); err != nil {
		return err
	}
	if err := g.SetKeybinding("information", gocui.KeyCtrlH, gocui.ModNone, helpToBodyView); err != nil {
		return err
	}
	/*if err := g.SetKeybinding("information", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("information", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}*/
	if err := g.SetKeybinding("body", gocui.KeyCtrlH, gocui.ModNone, bodyToHelpView); err != nil {
		return err
	}
	return nil
}
func saveDeleteView(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("savename"); err != nil {
		return err
	}
	s := v.Buffer()
	re := regexp.MustCompile(`\s+`)
	file = re.ReplaceAllString(s, "")
	v, err := g.SetCurrentView("body")
	if err != nil {
		return err
	}
	if err := saveMain(g, v); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("body", 0, 3, 2*maxX/3, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Dockerfile"
		v.Editable = true
		v.Wrap = true
		v.Editor = DockerfileEditor{}
	}
	if err := g.DeleteView("title"); err != nil {
		return err
	}
	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
	}
	return nil
}

func helpToBodyView(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("body"); err == nil {
		g.Mouse = false
		g.Cursor = true
		return nil
	}
	return nil
}

func bodyToHelpView(g *gocui.Gui, v *gocui.View) error {
	if viewExist(g, "information") {
		if _, err := g.SetCurrentView("information"); err == nil {
			v.Editable = false
			g.Mouse = false
			g.Cursor = false
			return err
		}
	} else {
		if v, err := g.SetCurrentView("help"); err == nil {
			v.Editable = false
			g.Mouse = false
			g.Cursor = false
			return err
		}
	}
	return nil
}

func runGocui(cmd *cobra.Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	file = dir + "/Dockerfile"
	if len(args) > 0 {
		file = dir + "/" + args[0]
	}
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = false
	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func saveView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("savename", 0, maxY/2-1, maxX, maxY/2+1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, file)
		v.Editable = true
	}
	if _, err := g.SetCurrentView("savename"); err != nil {
		return err
	}
	return nil
}

func main() {
	initializeMap()
	var rootCmd = &cobra.Command{
		Use:   "nanodock [file]",
		Short: "Terminal based Dockerfile Editor",
		Long:  "A Terminal based Dockerfile Editor",
		//Args:  cobra.MinimumNArgs(1),
		Run: runGocui,
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

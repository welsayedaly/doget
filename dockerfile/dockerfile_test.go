package dockerfile

import (
	"reflect"
	"strings"
	"testing"
)

type field interface{}

func assertEqual(expect, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Items not equal:\nexpected %q\nhave     %q\n", expect, actual)
	}
}

func assertParsed(expect interface{}, fieldOf func(d Dockerfile) field, input string, t *testing.T) {
	var fixture Dockerfile

	err := Parse(strings.NewReader(input), &fixture)
	if err != nil {
		t.Error("Could not parse " + err.Error())
	}

	assertEqual(expect, fieldOf(fixture), t)
}

func Test_parsing(t *testing.T) {
	file := `
FROM debian:jessie

# Fetch dependencies
# See https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/
RUN apt-get update && apt-get install -y \
    aufs-tools \
    automake \
    build-essential \
    curl \
    dpkg-sig \
    libcap-dev \
    libsqlite3-dev \
    mercurial \
    reprepro \
    ruby1.9.1 \
    ruby1.9.1-dev \
    s3cmd=1.1.* \
 && rm -rf /var/lib/apt/lists/*

CMD /bin/bash
	`

	assertParsed(4, func(d Dockerfile) field { return len(d.Statements) }, file, t)
}

func Test_parsing_comment(t *testing.T) {
	assertParsed("Hello", func(d Dockerfile) field { return d.Statements[0].(*Comment).Lines }, "# Hello", t)
}

func Test_parsing_multiline_comment(t *testing.T) {
	assertParsed("Hello\nWorld", func(d Dockerfile) field { return d.Statements[0].(*Comment).Lines }, "# Hello\n# World", t)
}

func Test_parsing_multiline_comment_with_empty_line(t *testing.T) {
	assertParsed("Hello\n\nWorld", func(d Dockerfile) field { return d.Statements[0].(*Comment).Lines }, "# Hello\n#\n# World", t)
}

func Test_parsing_from(t *testing.T) {
	assertParsed("scratch", func(d Dockerfile) field { return d.From.Image }, "FROM scratch", t)
}

func Test_parsing_from_statement(t *testing.T) {
	assertParsed("scratch", func(d Dockerfile) field { return d.Statements[0].(*From).Image }, "FROM scratch", t)
}

func Test_parsing_maintainer(t *testing.T) {
	assertParsed("Test", func(d Dockerfile) field { return d.Statements[0].(*Maintainer).Name }, "MAINTAINER Test", t)
}

func Test_parsing_run(t *testing.T) {
	assertParsed("apt-get update", func(d Dockerfile) field { return d.Statements[0].(*Run).Command }, "RUN apt-get update", t)
}

func Test_parsing_cmd(t *testing.T) {
	assertParsed("[\"/bin/bash\"]", func(d Dockerfile) field { return d.Statements[0].(*Cmd).CmdLine }, "CMD [\"/bin/bash\"]", t)
}

func Test_parsing_label(t *testing.T) {
	assertParsed("version=one", func(d Dockerfile) field { return d.Statements[0].(*Label).Pairs }, "LABEL version=one", t)
}

func Test_parsing_mutiline_label(t *testing.T) {
	assertParsed("label1=one\n      label2=two", func(d Dockerfile) field { return d.Statements[0].(*Label).Pairs }, "LABEL label1=one\\\n      label2=two", t)
}

func Test_parsing_expose(t *testing.T) {
	assertParsed("80", func(d Dockerfile) field { return d.Statements[0].(*Expose).Ports }, "EXPOSE 80", t)
}

func Test_parsing_env(t *testing.T) {
	assertParsed("PROFILE=dev", func(d Dockerfile) field { return d.Statements[0].(*Env).Pairs }, "ENV PROFILE=dev", t)
}

func Test_parsing_add(t *testing.T) {
	assertParsed("src /app", func(d Dockerfile) field { return d.Statements[0].(*Add).Paths }, "ADD src /app", t)
}

func Test_parsing_copy(t *testing.T) {
	assertParsed("src /app", func(d Dockerfile) field { return d.Statements[0].(*Copy).Paths }, "COPY src /app", t)
}

func Test_parsing_volume(t *testing.T) {
	assertParsed("/data", func(d Dockerfile) field { return d.Statements[0].(*Volume).Names }, "VOLUME /data", t)
}

func Test_parsing_user(t *testing.T) {
	assertParsed("app", func(d Dockerfile) field { return d.Statements[0].(*User).Name }, "USER app", t)
}

func Test_parsing_workdir(t *testing.T) {
	assertParsed("/data", func(d Dockerfile) field { return d.Statements[0].(*Workdir).Path }, "WORKDIR /data", t)
}

func Test_parsing_arg(t *testing.T) {
	assertParsed("name", func(d Dockerfile) field { return d.Statements[0].(*Arg).Name }, "ARG name", t)
}

func Test_parsing_onbuild(t *testing.T) {
	assertParsed("ADD . /app/src", func(d Dockerfile) field { return d.Statements[0].(*Onbuild).Instruction }, "ONBUILD ADD . /app/src", t)
}

func Test_parsing_stopsignal(t *testing.T) {
	assertParsed("SIGKILL", func(d Dockerfile) field { return d.Statements[0].(*Stopsignal).Signal }, "STOPSIGNAL SIGKILL", t)
}

func Test_parsing_healthcheck(t *testing.T) {
	assertParsed("NONE", func(d Dockerfile) field { return d.Statements[0].(*Healthcheck).Command }, "HEALTHCHECK NONE", t)
}

func Test_parsing_shell(t *testing.T) {
	assertParsed("[\"powershell\", \"-command\"]", func(d Dockerfile) field { return d.Statements[0].(*Shell).CmdLine }, "SHELL [\"powershell\", \"-command\"]", t)
}
FROM debian:bullseye-slim

RUN apt update && apt install -y apt-transport-https curl gnupg
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor > bazel-archive-keyring.gpg
RUN mv bazel-archive-keyring.gpg /usr/share/keyrings
RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/bazel-archive-keyring.gpg] https://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list
RUN apt update && apt install -y m4 clang-11 libssl-dev libcurl4-openssl-dev libarchive-dev libsqlite3-dev libmicrohttpd-dev pkg-config git bazel-6.4.0

RUN update-alternatives --install /usr/bin/clang clang /usr/bin/clang-11 380 && \
    update-alternatives --install /usr/bin/llc llc /usr/bin/llc-11 380 && \
    update-alternatives --install /usr/bin/llvm-strip llvm-strip /usr/bin/llvm-strip-11 380

WORKDIR /app/ebpf-transport-monitoring

COPY ebpf-transport-monitoring /app/ebpf-transport-monitoring

RUN bazel-6.4.0 build :lightfoot

ENTRYPOINT ["./bazel-bin/lightfoot", "--stdin_eof=false"]
#CMD ["tail", "-f", "/dev/null"]

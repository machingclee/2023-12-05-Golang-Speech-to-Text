FROM bitnami/golang:1.21-debian-11
RUN go env -w GO111MODULE=on
RUN go env -w GOSUMDB=sum.golang.org
#speechsdk start

# RUN mkdir -p /app/main
# WORKDIR /app/main
# COPY . /app/main
# RUN go mod tidy
ENV GOSUMDB="sum.golang.org"
ENV SPEECHSDK_ROOT="$HOME/speechsdk"

RUN apt-get update && apt-get install -y build-essential libssl-dev ca-certificates libasound2 wget ffmpeg \
    && mkdir -p "$SPEECHSDK_ROOT" \
    && wget -O SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary \
    && tar --strip 1 -xzf SpeechSDK-Linux.tar.gz -C "$SPEECHSDK_ROOT" \
    && ls -l "$SPEECHSDK_ROOT" \
    && rm SpeechSDK-Linux.tar.gz


ENV CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
ENV CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"
#speechsdk end
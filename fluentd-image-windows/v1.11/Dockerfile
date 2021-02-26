ARG SERVERCORE_VERSION
FROM mcr.microsoft.com/windows/servercore:${SERVERCORE_VERSION}
LABEL Description="Fluentd docker image" Vendor="Banzai Cloud" Version="1.11.5"

SHELL ["powershell", "-NoLogo", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]

RUN Set-ExecutionPolicy Bypass -Scope Process -Force; \
    iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

RUN choco install -y ruby --version 2.7.2.1 --params "'/InstallDir:C:\'"; \
    choco install -y msys2 --params "'/NoPath /NoUpdate /InstallDir:C:\ruby27\msys64'"

# ridk install 2 3; \
RUN refreshenv; \
'gem: --no-document' | Out-File -Encoding UTF8 -NoNewline -Force -FilePath 'C:\ProgramData\gemrc'; \
gem install oj -v 3.10.15; \
gem install http_parser.rb -v 0.5.3; \
gem install tzinfo -v 1.2.7; \
gem install json -v 2.3.1; \
gem install async-http -v 0.52.5; \
gem install ext_monitor -v 0.1.2; \
gem install fluentd -v 1.11.5; \
gem install prometheus-client -v 0.9.0; \
gem install bigdecimal -v 1.4.4; \
gem install webrick; \
gem install \
       specific_install \
       fluent-plugin-remote-syslog \
       fluent-plugin-webhdfs \
       fluent-plugin-elasticsearch \
       fluent-plugin-prometheus \
       fluent-plugin-s3 \
       fluent-plugin-rewrite-tag-filter \
       fluent-plugin-azure-storage-append-blob \
       fluent-plugin-oss \
       fluent-plugin-dedot_filter \
       fluent-plugin-sumologic_output \
       fluent-plugin-kafka \
       fluent-plugin-geoip \
       fluent-plugin-label-router \
       fluent-plugin-tag-normaliser \
       fluent-plugin-grafana-loki \
       fluent-plugin-concat \
       fluent-plugin-kinesis \
       fluent-plugin-parser-logfmt \
       fluent-plugin-detect-exceptions \
       fluent-plugin-multi-format-parser \
       fluent-plugin-record-modifier \
       fluent-plugin-newrelic \
       fluent-plugin-splunk-hec \
       elasticsearch-xpack \
       fluent-plugin-cloudwatch-logs \
       fluent-plugin-throttle \
       fluent-plugin-logdna \
       fluent-plugin-datadog \
       fluent-plugin-aws-elasticsearch-service \
       fluent-plugin-redis; \
gem specific_install -l https://github.com/tarokkk/fluent-plugin-logzio.git; \
gem specific_install -l https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424.git; \
gem specific_install -l https://github.com/banzaicloud/fluent-plugin-gcs.git; \
gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git -d fluent-plugin-kubernetes-metadata-filter; \
gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git -d fluent-plugin-enhance-k8s-metadata; \
gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git -d fluent-plugin-kubernetes-sumologic; \
Remove-Item -Force C:\ruby27\lib\ruby\gems\2.7.2\cache\*.gem; \
Remove-Item -Recurse -Force C:\ProgramData\chocolatey; \
mkdir C:\fluentd\log C:\fluentd\etc C:\buffers

COPY fluent.conf 'C:\etc\fluent\fluent.conf'
ENV FLUENTD_CONF="fluent.conf"
EXPOSE 24224 5140
CMD ["powershell", "-command", "fluentd"]

FROM circleci/golang

RUN sudo mkdir /.config
RUN sudo chown circleci:circleci /.config
ENV HOME /
ENV CLOUDSDK_PYTHON_SITEPACKAGES 1
RUN wget https://dl.google.com/dl/cloudsdk/channels/rapid/google-cloud-sdk.zip && unzip google-cloud-sdk.zip && rm google-cloud-sdk.zip

ENV CLOUDSDK_COMPONENT_MANAGER_FIXED_SDK_VERSION 0.9.82
RUN sudo google-cloud-sdk/install.sh --usage-reporting=true --path-update=true --bash-completion=true --rc-path=/.bashrc --additional-components app-engine-java app-engine-python app-engine-go kubectl alpha beta
RUN sudo google-cloud-sdk/bin/gcloud config set --installation component_manager/disable_update_check true || google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check true
RUN sed -i -- 's/\"disable_updater\": false/\"disable_updater\": true/g' /google-cloud-sdk/lib/googlecloudsdk/core/config.json || echo nope, too soon for this
ENV CLOUD_SDK_PATH /google-cloud-sdk/bin
ENV PATH $CLOUD_SDK_PATH:$PATH
VOLUME ["/.config"]
CMD bash
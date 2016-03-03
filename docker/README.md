### Usage

```
$ docker run --rm -i -t \
    -e WORKDIR=repo \
    -e REPOSITORY=<REPOSITORY_URL> \
    -e UNITY3D2PNG_URL=http://$(docker-machine ip default):19300 \
    -e PATH_TEMPLATE="path/to/{{with .req_param1}}{{index . 0}}{{else}}default{{end}}/{{with .req_param2}}{{index . 0}}{{else}}default{{end}}/{{index .file 0}}" \
    -v $HOME/.ssh/id_rsa:/root/.ssh/id_rsa \
    -p 19301:19300 \
    <DOCKER_IMAGE>
```
```
$ curl "http://$(docker-machine ip default):19301?branch=master&file=hoge.unity3d" -s > /tmp/output.png
$ curl "http://$(docker-machine ip default):19301?branch=master&file=hoge.unity3d&req_param=1&req_param2" -s > /tmp/output.png
```

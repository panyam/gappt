# Go App Template

A minimal template app for go backends that require:

1. Protos/GRPC services
2. Fronted by a gateway service
3. And powered by oauth
4. Basic frontend based on go-templates - this can be customized in the future.
5. Common docker compose manifests for packaging for development.
6. Optional k8s configs if needed in the future for testing against cluster deployments

Other backend choices (like datastores) are upto the app/service dev, eg:

1. Which services are to be added.
2. Which backends are to be used (for storage, etc)
3. How to deploy them to specific hosting providers (eg appengine, heroku etc)
4. Selecting frontend frameworks.

## Requirements

1. Basic/Standard Go Tooling:

* Go
* Air (for fast reloads)
* Protobuf
* GRPC
* Buf (for generating artificates from grpc protos)

2. Python (for repo setup)

## Getting Started

1. Clone this Repo
2. Setup your variables in VARS.yaml.  This will let you configure things like repo names etc.   See the comments in VARS.yaml for how to customize your target repo.
3. Run the bootstrapper:

```
python bootstrap.py VARS.yaml <output_dir>
```

The "output_dir" *can* exist any where.   Only conditions are that it cannot be "this" folder and that the files being generated must not exist otherwise an error is thrown.   To override this behavior (and forcefully replace existing files) use the `--override` flag.

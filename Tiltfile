# Welcome to Tilt!
#   To get you started as quickly as possible, we have created a
#   starter Tiltfile for you.
#
#   Uncomment, modify, and delete any commands as needed for your
#   project's configuration.


# Output diagnostic messages
#   You can print log messages, warnings, and fatal errors, which will
#   appear in the (Tiltfile) resource in the web UI. Tiltfiles support
#   multiline strings and common string operations such as formatting.
#
#   More info: https://docs.tilt.dev/api.html#api.warn
print("""
-----------------------------------------------------------------
✨ Hello Tilt! This appears in the (Tiltfile) pane whenever Tilt
   evaluates this file.
-----------------------------------------------------------------
""".strip())
warn('ℹ️ Open {tiltfile_path} in your favorite editor to get started.'.format(
    tiltfile_path=config.main_path))


# Build Docker image
#   Tilt will automatically associate image builds with the resource(s)
#   that reference them (e.g. via Kubernetes or Docker Compose YAML).
#
#   More info: https://docs.tilt.dev/api.html#api.docker_build
#
# docker_build('registry.example.com/my-image',
#              context='.',
#              # (Optional) Use a custom Dockerfile path
#              dockerfile='./deploy/app.dockerfile',
#              # (Optional) Filter the paths used in the build
#              only=['./app'],
#              # (Recommended) Updating a running container in-place
#              # https://docs.tilt.dev/live_update_reference.html
#              live_update=[
#                 # Sync files from host to container
#                 sync('./app', '/src/'),
#                 # Execute commands inside the container when certain
#                 # paths change
#                 run('/src/codegen.sh', trigger=['./app/api'])
#              ]
#)
# Build Docker image for service 1
docker_build('err/backend-service1',
context="./backend/service",
dockerfile='service1-dockerfile',
live_update=[
    run(
        'CGO_ENABLED=0 GOOS=linux go build -o ./build ./cmd/service1/main.go',
        trigger=[
            './backend/service/cmd/service1/*',
            './backend/service/shared/*'
        ]
    )
]
)

docker_build('err/backend-service2',
context="./backend/service",
dockerfile='service2-dockerfile',
live_update=[
    # run(
    #     'CGO_ENABLED=0 GOOS=linux go build -o ./build ./cmd/service2/main.go',
    #     trigger=[
    #         './backend/service/cmd/service2/*',
    #         './backend/service/shared/*'
    #     ]
    # ),
]
)

docker_build('err/github-app',
context="./backend/service",
dockerfile='ghapp.dockerfile',
live_update=[
    run(
        'CGO_ENABLED=0 GOOS=linux go build -o ./build ./cmd/github-app/main.go',
        trigger=[
            './backend/service/cmd/github-app/*',
            './backend/service/shared/*'
        ]
    ),
]
)

docker_build('err/frontend',
context="./frontend",
dockerfile='frontend.dockerfile',
)

docker_build(
        ref = "solana-contract",
        context = "./backend/identity",
        dockerfile = "solana.dockerfile",
        target = "builder",
        build_args = {"BRIDGE_ADDRESS": "Bridge1p5gheXUvJ6jGWGeCsgPKgnE3YgdGKRVCMY9o"}
    )

k8s_yaml([
    'k8s/global.configMap.yaml',
    'k8s/secret.vault.yaml', 
    'k8s/service1.postgres.yaml',
    'k8s/service1.deployment.yaml',
    'k8s/service2.deployment.yaml',
    'k8s/ghapp.deployment.yaml',
    'k8s/ghapp.postgres.yaml',
    'k8s/frontend.deployment.yaml',
    'k8s/solana.devnet.yaml',
])

# Apply Kubernetes manifests
#   Tilt will build & push any necessary images, re-deploying your
#   resources as they change.
#
#   More info: https://docs.tilt.dev/api.html#api.k8s_yaml
#
# k8s_yaml(['k8s/deployment.yaml', 'k8s/service.yaml'])


# Customize a Kubernetes resource
#   By default, Kubernetes resource names are automatically assigned
#   based on objects in the YAML manifests, e.g. Deployment name.
#
#   Tilt strives for sane defaults, so calling k8s_resource is
#   optional, and you only need to pass the arguments you want to
#   override.
#
#   More info: https://docs.tilt.dev/api.html#api.k8s_resource
#
# k8s_resource('my-deployment',
#              # map one or more local ports to ports on your Pod
#              port_forwards=['5000:8080'],
#              # change whether the resource is started by default
#              auto_init=False,
#              # control whether the resource automatically updates
#              trigger_mode=TRIGGER_MODE_MANUAL
# )


k8s_resource("github-app",resource_deps=['global-configmap','secret-vault','ghapp-psql'],port_forwards=["30005:8080"])
k8s_resource('microservice2',resource_deps=['microservice1'],port_forwards=["30004:8080"])
k8s_resource('microservice1',resource_deps=['postgres'],port_forwards=['30002:1122',"30003:8080"])
k8s_resource('frontend',port_forwards=["33030:3030"])
k8s_resource('postgres',
resource_deps=['global-configmap','secret-vault'],
                port_forwards=['30001:5432'],
    )

k8s_resource('ghapp-psql',
resource_deps=['global-configmap','secret-vault'],
                port_forwards=['30006:5432'],
    )
local_resource('migrate_service1',
               cmd='just migrate_up_s1 pwd',
            resource_deps=['postgres']
)
local_resource('migrate_ghapp',
               cmd='just migrate_up_ghapp pwd',
            resource_deps=['ghapp-psql']
)
k8s_resource(
        "solana-devnet",
        port_forwards = [
            port_forward(8899, name = "Solana RPC [:8899]", host = "localhost"),
            port_forward(8900, name = "Solana WS [:8900]", host = "localhost"),
        ],
        labels = ["solana"],
    )
## Install kafka helm 
load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('bitnami', 'https://charts.bitnami.com/bitnami')
helm_resource('kafka', 'bitnami/kafka')


# Run local commands
#   Local commands can be helpful for one-time tasks like installing
#   project prerequisites. They can also manage long-lived processes
#   for non-containerized services or dependencies.
#
#   More info: https://docs.tilt.dev/local_resource.html
#



# Extensions are open-source, pre-packaged functions that extend Tilt
#
#   More info: https://github.com/tilt-dev/tilt-extensions
#
load('ext://git_resource', 'git_checkout')


# Organize logic into functions
#   Tiltfiles are written in Starlark, a Python-inspired language, so
#   you can use functions, conditionals, loops, and more.
#
#   More info: https://docs.tilt.dev/tiltfile_concepts.html
#
def tilt_demo():
    # Tilt provides many useful portable built-ins
    # https://docs.tilt.dev/api.html#modules.os.path.exists
    if os.path.exists('tilt-avatars/Tiltfile'):
        # It's possible to load other Tiltfiles to further organize
        # your logic in large projects
        # https://docs.tilt.dev/multiple_repos.html
        load_dynamic('tilt-avatars/Tiltfile')
    watch_file('tilt-avatars/Tiltfile')
    git_checkout('https://github.com/tilt-dev/tilt-avatars.git',
                 checkout_dir='tilt-avatars')


# Edit your Tiltfile without restarting Tilt
#   While running `tilt up`, Tilt watches the Tiltfile on disk and
#   automatically re-evaluates it on change.
#
#   To see it in action, try uncommenting the following line with
#   Tilt running.
# tilt_demo()

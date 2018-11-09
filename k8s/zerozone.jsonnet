local kube = import 'kube.libsonnet';

{
  namespace:: { metadata+: { namespace: 'zerozone' } },
  ns: kube.Namespace($.namespace.metadata.namespace),

  svc: kube.Service('zerozone') + $.namespace {
    target_pod: $.zerozone.spec.template,
    port: 53,
    spec+: {
      type: 'LoadBalancer',

      // kube.Service failes to copy the protocol from the target_pod
      local sport = super.ports[0],
      ports: [
        sport {
          protocol: 'UDP',
        },
      ],
    },
  },

  cfg: kube.ConfigMap('zerozone') + $.namespace {
    data: {
      Corefile: |||
        0zone.mkm.pub:8053 {
            zerozone ipfs:5001
            prometheus localhost:9253
            errors
            log
        }
      |||,
    },
  },

  zerozone: kube.Deployment('zerozone') + $.namespace {
    spec+: {
      template+: {
        spec+: {
          default_container: 'zerozone_server',
          containers_+: {
            debug: kube.Container('debug') {
              image: 'ubuntu',
              args: ['/bin/sleep', '10000000'],
              volumeMounts_+: {
                cfg: {
                  mountPath: '/Corefile',
                  subPath: 'Corefile',
                },
              },
              resources+: {
                requests+: { memory: '10Mi' },
              },
            },

            zerozone_server: kube.Container('zerozone-server') {
              image: 'mkmik/zerozone-server@sha256:24d8f8935c17f8b2f3283b263a666d8adfeab5ffe395de31ffef73a06c1065ee',
              ports_+: {
                dns: { containerPort: 8053, protocol: 'UDP' },
              },
              volumeMounts_+: {
                cfg: {
                  mountPath: '/Corefile',
                  subPath: 'Corefile',
                },
              },
              resources+: {
                requests+: { memory: '10Mi' },
              },
            },
          },
          volumes_: +{
            cfg: {
              configMap: { name: 'zerozone' },
            },
          },
        },

      },
    },
  },

  ipfsSvc: kube.Service('ipfs') + $.namespace {
    target_pod: $.ipfs.spec.template,
    spec+: {
      ports: [
        { name: 'api', port: 5001 },
      ],
    },
  },

  // zerozone dns server will talk to this local ipfs node set in order to fetch zones
  ipfs: kube.Deployment('ipfs') + $.namespace {
    spec+: {
      template+: {
        spec+: {
          containers_+: {
            go_ipfs: kube.Container('go-ipfs') {
              image: 'ipfs/go-ipfs:v0.4.18',
              args: ['daemon', '--enable-namesys-pubsub', '--enable-pubsub-experiment'],
              ports_+: {
                api: { containerPort: 5001 },
              },
              volumeMounts_+: {
                data: { mountPath: '/data/ipfs' },
              },
              resources+: {
                requests+: { memory: '10Mi' },
              },
            },
          },
          volumes_: +{
            data: {
              emptyDir: {},
            },
          },
        },
      },
    },
  },
}

local kube = import 'kube.libsonnet';

{
  namespace:: { metadata+: { namespace: 'zerozone' } },
  ns: kube.Namespace($.namespace.metadata.namespace),

  svc: kube.Service('zerozone') + $.namespace {
    target_pod: $.zerozone_server.spec.template,
    port: 53,
    spec+: {
      type: 'LoadBalancer',

      // kube.Service failes to copy the protocol from the target_pod
      local sport = super.ports[0],
      ports: [sport {
        protocol: 'UDP',
      }],
    },
  },

  zerozone_server: kube.Deployment('zerozone-server') + $.namespace {
    spec+: {
      template+: {
        spec+: {
          containers_+: {
            zerozone_server: kube.Container('zerozone-server') {
              image: 'mkmik/zerozone-server@sha256:24d8f8935c17f8b2f3283b263a666d8adfeab5ffe395de31ffef73a06c1065ee',
              ports_+: {
                dns: { containerPort: 8053, protocol: 'UDP' },
              },
              resources+: {
                requests+: { memory: '10Mi' },
              },
            },
          },
        },
      },
    },
  },
}
